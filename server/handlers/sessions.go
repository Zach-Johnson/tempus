package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	pb "github.com/Zach-Johnson/tempus/proto/api/v1/tempus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// PracticeSessionHandler implements the PracticeSessionService gRPC service
type PracticeSessionHandler struct {
	pb.UnimplementedPracticeSessionServiceServer
	db *sql.DB
}

// NewPracticeSessionHandler creates a new PracticeSessionHandler
func NewPracticeSessionHandler(db *sql.DB) *PracticeSessionHandler {
	return &PracticeSessionHandler{db: db}
}

// CreatePracticeSession creates a new practice session
func (h *PracticeSessionHandler) CreatePracticeSession(ctx context.Context, req *pb.CreatePracticeSessionRequest) (*pb.PracticeSession, error) {
	// Validate request
	if req.StartTime == nil {
		return nil, status.Error(codes.InvalidArgument, "start time is required")
	}
	if req.EndTime == nil {
		return nil, status.Error(codes.InvalidArgument, "end time is required")
	}

	startTime := req.StartTime.AsTime()
	endTime := req.EndTime.AsTime()

	if startTime.After(endTime) {
		return nil, status.Error(codes.InvalidArgument, "start time cannot be after end time")
	}

	// Start a transaction
	tx, err := h.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to start transaction: %v", err)
	}
	defer tx.Rollback() // Rollback if not committed

	// Insert the practice session
	result, err := tx.ExecContext(
		ctx,
		"INSERT INTO practice_sessions (start_time, end_time, notes) VALUES (?, ?, ?)",
		startTime, endTime, req.Notes,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create practice session: %v", err)
	}

	// Get the session ID
	sessionID, err := result.LastInsertId()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get session ID: %v", err)
	}

	// Get the session creation and update times
	var createdAt, updatedAt time.Time
	err = tx.QueryRowContext(
		ctx,
		"SELECT created_at, updated_at FROM practice_sessions WHERE id = ?",
		sessionID,
	).Scan(&createdAt, &updatedAt)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get session times: %v", err)
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to commit transaction: %v", err)
	}

	// Return the created session
	return &pb.PracticeSession{
		Id:        int32(sessionID),
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		Notes:     req.Notes,
		CreatedAt: timestamppb.New(createdAt),
		UpdatedAt: timestamppb.New(updatedAt),
	}, nil
}

// GetPracticeSession retrieves a practice session by ID
func (h *PracticeSessionHandler) GetPracticeSession(ctx context.Context, req *pb.GetPracticeSessionRequest) (*pb.PracticeSession, error) {
	if req.Id <= 0 {
		return nil, status.Error(codes.InvalidArgument, "invalid session ID")
	}

	// Start a transaction to ensure consistency across queries
	tx, err := h.db.BeginTx(ctx, &sql.TxOptions{ReadOnly: true})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to start transaction: %v", err)
	}
	defer tx.Rollback() // Rollback if not committed

	// Query the session
	var session pb.PracticeSession
	var startTime, endTime, createdAt, updatedAt time.Time

	err = tx.QueryRowContext(
		ctx,
		"SELECT id, start_time, end_time, notes, created_at, updated_at FROM practice_sessions WHERE id = ?",
		req.Id,
	).Scan(&session.Id, &startTime, &endTime, &session.Notes, &createdAt, &updatedAt)

	if err == sql.ErrNoRows {
		return nil, status.Errorf(codes.NotFound, "practice session with ID %d not found", req.Id)
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to retrieve practice session: %v", err)
	}

	session.StartTime = timestamppb.New(startTime)
	session.EndTime = timestamppb.New(endTime)
	session.CreatedAt = timestamppb.New(createdAt)
	session.UpdatedAt = timestamppb.New(updatedAt)

	// Query exercise history
	exerciseRows, err := tx.QueryContext(
		ctx,
		`SELECT id, exercise_id, start_time, end_time, bpms, time_signature, notes
         FROM exercise_history
         WHERE session_id = ?
         ORDER BY start_time`,
		req.Id,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to retrieve session exercises: %v", err)
	}
	defer exerciseRows.Close()

	// Parse exercises
	var exercises []*pb.ExerciseHistory
	exerciseIDMap := make(map[int32]int) // Maps exercise ID to index in exercises slice

	for exerciseRows.Next() {
		var sessionExercise pb.ExerciseHistory
		var exerciseId int32
		var startTime, endTime time.Time
		var bpmJSON string

		err := exerciseRows.Scan(
			&sessionExercise.Id,
			&exerciseId,
			&startTime,
			&endTime,
			&bpmJSON,
			&sessionExercise.TimeSignature,
			&sessionExercise.Notes,
		)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to parse session exercise: %v", err)
		}

		if bpmJSON != "" {
			var bpms []int32
			if err := json.Unmarshal([]byte(bpmJSON), &bpms); err != nil {
				return nil, status.Errorf(codes.Internal, "failed to unmarshal BPM values: %v", err)
			}
			sessionExercise.Bpms = bpms
		}

		sessionExercise.SessionId = req.Id
		sessionExercise.ExerciseId = exerciseId
		sessionExercise.StartTime = timestamppb.New(startTime)
		sessionExercise.EndTime = timestamppb.New(endTime)

		// Fetch exercise details
		exercise, err := h.getExerciseDetails(ctx, tx, exerciseId)
		if err != nil {
			return nil, err
		}
		sessionExercise.Exercise = exercise

		exercises = append(exercises, &sessionExercise)
		exerciseIDMap[exerciseId] = len(exercises) - 1
	}

	if err = exerciseRows.Err(); err != nil {
		return nil, status.Errorf(codes.Internal, "error reading session exercises: %v", err)
	}

	session.Exercises = exercises

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to commit transaction: %v", err)
	}

	return &session, nil
}

// ListPracticeSessions lists practice sessions with optional filtering and pagination
func (h *PracticeSessionHandler) ListPracticeSessions(ctx context.Context, req *pb.ListPracticeSessionsRequest) (*pb.ListPracticeSessionsResponse, error) {
	pageSize := int(req.PageSize)
	if pageSize <= 0 {
		pageSize = 50 // Default page size
	}

	// Parse the page token, which is just an offset
	offset := 0
	if req.PageToken != "" {
		var err error
		offset, err = strconv.Atoi(req.PageToken)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "invalid page token")
		}
	}

	// Build the query based on filters
	baseQuery := `
        SELECT id, start_time, end_time, notes, created_at, updated_at
        FROM practice_sessions
    `
	countQuery := `
        SELECT COUNT(*)
        FROM practice_sessions
    `

	var whereClause string
	var queryParams []interface{}

	// Add filter by date range if provided
	if req.StartDate != nil && req.EndDate != nil {
		whereClause = " WHERE start_time >= ? AND end_time <= ?"
		queryParams = append(queryParams, req.StartDate.AsTime(), req.EndDate.AsTime())
	} else if req.StartDate != nil {
		whereClause = " WHERE start_time >= ?"
		queryParams = append(queryParams, req.StartDate.AsTime())
	} else if req.EndDate != nil {
		whereClause = " WHERE end_time <= ?"
		queryParams = append(queryParams, req.EndDate.AsTime())
	}

	// Add filter by exercise if provided
	if req.ExerciseId > 0 {
		if whereClause == "" {
			whereClause = " WHERE"
		} else {
			whereClause += " AND"
		}
		whereClause += " id IN (SELECT session_id FROM exercise_history WHERE exercise_id = ?)"
		queryParams = append(queryParams, req.ExerciseId)
	}

	// Add order by, limit, and offset
	fullQuery := baseQuery + whereClause + " ORDER BY start_time DESC LIMIT ? OFFSET ?"
	queryParams = append(queryParams, pageSize+1, offset) // Query one more to check if there are more pages

	// Query total count
	var totalCount int32
	countQueryParams := make([]interface{}, len(queryParams)-2) // Exclude limit and offset
	copy(countQueryParams, queryParams[:len(queryParams)-2])

	err := h.db.QueryRowContext(ctx, countQuery+whereClause, countQueryParams...).Scan(&totalCount)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to count practice sessions: %v", err)
	}

	// Query sessions
	rows, err := h.db.QueryContext(ctx, fullQuery, queryParams...)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list practice sessions: %v", err)
	}
	defer rows.Close()

	// Parse the results
	sessions := make([]*pb.PracticeSession, 0, pageSize)
	count := 0
	hasMorePages := false

	for rows.Next() {
		if count >= pageSize {
			hasMorePages = true
			break
		}

		var session pb.PracticeSession
		var startTime, endTime, createdAt, updatedAt time.Time

		err := rows.Scan(
			&session.Id, &startTime, &endTime, &session.Notes, &createdAt, &updatedAt,
		)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to parse practice session: %v", err)
		}

		session.StartTime = timestamppb.New(startTime)
		session.EndTime = timestamppb.New(endTime)
		session.CreatedAt = timestamppb.New(createdAt)
		session.UpdatedAt = timestamppb.New(updatedAt)
		session.Exercises = []*pb.ExerciseHistory{}

		sessions = append(sessions, &session)
		count++
	}

	if err = rows.Err(); err != nil {
		return nil, status.Errorf(codes.Internal, "error reading practice sessions: %v", err)
	}

	// Calculate next page token
	nextPageToken := ""
	if hasMorePages {
		nextPageToken = strconv.Itoa(offset + pageSize)
	}

	return &pb.ListPracticeSessionsResponse{
		Sessions:      sessions,
		NextPageToken: nextPageToken,
		TotalCount:    totalCount,
	}, nil
}

// UpdatePracticeSession updates a practice session
func (h *PracticeSessionHandler) UpdatePracticeSession(ctx context.Context, req *pb.UpdatePracticeSessionRequest) (*pb.PracticeSession, error) {
	if req.Id <= 0 {
		return nil, status.Error(codes.InvalidArgument, "invalid session ID")
	}

	if req.Session == nil {
		return nil, status.Error(codes.InvalidArgument, "session data is required")
	}

	// Check if the session exists
	var exists bool
	err := h.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM practice_sessions WHERE id = ?)", req.Id).Scan(&exists)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to check session existence: %v", err)
	}
	if !exists {
		return nil, status.Errorf(codes.NotFound, "session with ID %d not found", req.Id)
	}

	// Parse update mask
	updateStartTime := false
	updateEndTime := false
	updateNotes := false

	if req.UpdateMask == nil || len(req.UpdateMask.Paths) == 0 {
		// If no update mask is provided, update all fields
		updateStartTime = true
		updateEndTime = true
		updateNotes = true
	} else {
		for _, path := range req.UpdateMask.Paths {
			switch path {
			case "start_time":
				updateStartTime = true
			case "end_time":
				updateEndTime = true
			case "notes":
				updateNotes = true
			}
		}
	}

	// Validate times if updating
	if updateStartTime && updateEndTime {
		startTime := req.Session.StartTime.AsTime()
		endTime := req.Session.EndTime.AsTime()
		if startTime.After(endTime) {
			return nil, status.Error(codes.InvalidArgument, "start time cannot be after end time")
		}
	} else if updateStartTime {
		// If only updating start time, check against existing end time
		var endTime time.Time
		err := h.db.QueryRowContext(ctx, "SELECT end_time FROM practice_sessions WHERE id = ?", req.Id).Scan(&endTime)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to get existing end time: %v", err)
		}

		startTime := req.Session.StartTime.AsTime()
		if startTime.After(endTime) {
			return nil, status.Error(codes.InvalidArgument, "start time cannot be after end time")
		}
	} else if updateEndTime {
		// If only updating end time, check against existing start time
		var startTime time.Time
		err := h.db.QueryRowContext(ctx, "SELECT start_time FROM practice_sessions WHERE id = ?", req.Id).Scan(&startTime)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to get existing start time: %v", err)
		}

		endTime := req.Session.EndTime.AsTime()
		if startTime.After(endTime) {
			return nil, status.Error(codes.InvalidArgument, "start time cannot be after end time")
		}
	}

	// Build update SQL
	sql := "UPDATE practice_sessions SET"
	params := []interface{}{}
	first := true

	if updateStartTime {
		if req.Session.StartTime == nil {
			return nil, status.Error(codes.InvalidArgument, "start time cannot be empty")
		}
		sql += " start_time = ?"
		params = append(params, req.Session.StartTime.AsTime())
		first = false
	}

	if updateEndTime {
		if req.Session.EndTime == nil {
			return nil, status.Error(codes.InvalidArgument, "end time cannot be empty")
		}
		if !first {
			sql += ","
		}
		sql += " end_time = ?"
		params = append(params, req.Session.EndTime.AsTime())
		first = false
	}

	if updateNotes {
		if !first {
			sql += ","
		}
		sql += " notes = ?"
		params = append(params, req.Session.Notes)
	}

	if len(params) == 0 {
		return nil, status.Error(codes.InvalidArgument, "no fields to update")
	}

	sql += " WHERE id = ?"
	params = append(params, req.Id)

	// Execute the update
	_, err = h.db.ExecContext(ctx, sql, params...)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update practice session: %v", err)
	}

	// Retrieve the updated session
	return h.GetPracticeSession(ctx, &pb.GetPracticeSessionRequest{Id: req.Id})
}

// DeletePracticeSession deletes a practice session
func (h *PracticeSessionHandler) DeletePracticeSession(ctx context.Context, req *pb.DeletePracticeSessionRequest) (*emptypb.Empty, error) {
	if req.Id <= 0 {
		return nil, status.Error(codes.InvalidArgument, "invalid session ID")
	}

	// Check if the session exists
	var exists bool
	err := h.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM practice_sessions WHERE id = ?)", req.Id).Scan(&exists)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to check session existence: %v", err)
	}
	if !exists {
		return nil, status.Errorf(codes.NotFound, "session with ID %d not found", req.Id)
	}

	// Delete the session (associated exercises will be deleted by ON DELETE CASCADE)
	_, err = h.db.ExecContext(ctx, "DELETE FROM practice_sessions WHERE id = ?", req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete practice session: %v", err)
	}

	return &emptypb.Empty{}, nil
}

// GetPracticeStats returns statistics for practice sessions
func (h *PracticeSessionHandler) GetPracticeStats(ctx context.Context, req *pb.GetPracticeStatsRequest) (*pb.PracticeStats, error) {
	// Build query filters based on request parameters
	var whereClause string
	var queryParams []interface{}

	// Add filter by date range if provided
	if req.StartDate != nil && req.EndDate != nil {
		whereClause = " WHERE ps.start_time >= ? AND ps.end_time <= ?"
		queryParams = append(queryParams, req.StartDate.AsTime(), req.EndDate.AsTime())
	} else if req.StartDate != nil {
		whereClause = " WHERE ps.start_time >= ?"
		queryParams = append(queryParams, req.StartDate.AsTime())
	} else if req.EndDate != nil {
		whereClause = " WHERE ps.end_time <= ?"
		queryParams = append(queryParams, req.EndDate.AsTime())
	}

	// Add category filter if provided
	var categoryJoin, categoryFilter string
	if req.CategoryId > 0 {
		categoryJoin = " JOIN exercise_history eh ON eh.session_id = ps.id JOIN exercise_categories ec ON eh.exercise_id = ec.exercise_id"
		if whereClause == "" {
			categoryFilter = " WHERE ec.category_id = ?"
		} else {
			categoryFilter = " AND ec.category_id = ?"
		}
		queryParams = append(queryParams, req.CategoryId)
	}

	// Get total sessions and duration
	var totalSessions int32
	var totalDurationSeconds int32

	sessionQuery := "SELECT COUNT(DISTINCT ps.id), COALESCE(SUM(strftime('%s', ps.end_time) - strftime('%s', ps.start_time)), 0) FROM practice_sessions ps" + categoryJoin + whereClause + categoryFilter
	err := h.db.QueryRowContext(ctx, sessionQuery, queryParams...).Scan(&totalSessions, &totalDurationSeconds)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to retrieve session statistics: %v", err)
	}

	// Calculate average session duration
	var avgSessionDurationSeconds float64
	if totalSessions > 0 {
		avgSessionDurationSeconds = float64(totalDurationSeconds) / float64(totalSessions)
	}

	// Get exercise time distribution
	exerciseDistQuery := `
		SELECT 
			e.id, 
			e.name, 
			COALESCE(SUM(strftime('%s', eh.end_time) - strftime('%s', eh.start_time)), 0) as duration,
			ROUND(COALESCE(SUM(strftime('%s', eh.end_time) - strftime('%s', eh.start_time)), 0) * 100.0 / ?, 2) as percentage
		FROM 
			exercises e
		JOIN 
			exercise_history eh ON e.id = eh.exercise_id
		JOIN 
			practice_sessions ps ON eh.session_id = ps.id
	` + whereClause + `
		GROUP BY 
			e.id
		ORDER BY 
			duration DESC
		LIMIT 10
	`

	exerciseDistParams := append([]interface{}{totalDurationSeconds}, queryParams...)
	exerciseRows, err := h.db.QueryContext(ctx, exerciseDistQuery, exerciseDistParams...)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to calculate exercise distribution: %v", err)
	}
	defer exerciseRows.Close()

	var exerciseDistribution []*pb.ExerciseTimeDistribution
	for exerciseRows.Next() {
		var dist pb.ExerciseTimeDistribution
		var percentage float64

		if err := exerciseRows.Scan(&dist.ExerciseId, &dist.ExerciseName, &dist.DurationSeconds, &percentage); err != nil {
			return nil, status.Errorf(codes.Internal, "failed to parse exercise distribution: %v", err)
		}

		dist.Percentage = percentage
		exerciseDistribution = append(exerciseDistribution, &dist)
	}

	if err = exerciseRows.Err(); err != nil {
		return nil, status.Errorf(codes.Internal, "error reading exercise distribution: %v", err)
	}

	// Get overall practice frequency by day
	frequencyQuery := `
		SELECT 
			date(ps.start_time) as practice_date,
			SUM(strftime('%s', ps.end_time) - strftime('%s', ps.start_time)) as duration
		FROM 
			practice_sessions ps
	` + whereClause + `
		GROUP BY 
			practice_date
		ORDER BY 
			practice_date ASC
	`

	frequencyRows, err := h.db.QueryContext(ctx, frequencyQuery, queryParams...)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to calculate practice frequency: %v", err)
	}
	defer frequencyRows.Close()

	var practiceFrequency []*pb.PracticeTimePoint
	for frequencyRows.Next() {
		var dateStr string
		var durationSeconds int32
		var practiceDate time.Time

		if err := frequencyRows.Scan(&dateStr, &durationSeconds); err != nil {
			return nil, status.Errorf(codes.Internal, "failed to parse practice frequency: %v", err)
		}

		// Parse the date string (format: YYYY-MM-DD)
		practiceDate, err = time.Parse("2006-01-02", dateStr)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to parse date: %v", err)
		}

		point := &pb.PracticeTimePoint{
			Date:            timestamppb.New(practiceDate),
			DurationSeconds: durationSeconds,
		}
		practiceFrequency = append(practiceFrequency, point)
	}

	if err = frequencyRows.Err(); err != nil {
		return nil, status.Errorf(codes.Internal, "error reading practice frequency: %v", err)
	}

	// Get category time distribution with daily breakdown
	categoryDistQuery := `
		SELECT 
			c.id, 
			c.name, 
			COALESCE(SUM(strftime('%s', eh.end_time) - strftime('%s', eh.start_time)), 0) as duration,
			ROUND(COALESCE(SUM(strftime('%s', eh.end_time) - strftime('%s', eh.start_time)), 0) * 100.0 / ?, 2) as percentage
		FROM 
			categories c
		JOIN 
			exercise_categories ec ON c.id = ec.category_id
		JOIN 
			exercise_history eh ON ec.exercise_id = eh.exercise_id
		JOIN 
			practice_sessions ps ON eh.session_id = ps.id
	` + whereClause + `
		GROUP BY 
			c.id
		ORDER BY 
			duration DESC
	`

	categoryDistParams := append([]interface{}{totalDurationSeconds}, queryParams...)
	categoryRows, err := h.db.QueryContext(ctx, categoryDistQuery, categoryDistParams...)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to calculate category distribution: %v", err)
	}
	defer categoryRows.Close()

	var categoryDistribution []*pb.CategoryTimeDistribution

	for categoryRows.Next() {
		var dist pb.CategoryTimeDistribution
		var percentage float64

		if err := categoryRows.Scan(&dist.CategoryId, &dist.CategoryName, &dist.DurationSeconds, &percentage); err != nil {
			return nil, status.Errorf(codes.Internal, "failed to parse category distribution: %v", err)
		}

		dist.Percentage = percentage

		// For each category, get daily practice times
		categoryFrequency, err := h.getCategoryDailyPractice(ctx, dist.CategoryId, whereClause, queryParams)
		if err != nil {
			return nil, err
		}

		dist.PracticeFrequency = categoryFrequency
		categoryDistribution = append(categoryDistribution, &dist)
	}

	if err = categoryRows.Err(); err != nil {
		return nil, status.Errorf(codes.Internal, "error reading category distribution: %v", err)
	}

	// Return the statistics
	return &pb.PracticeStats{
		TotalSessions:             totalSessions,
		TotalDurationSeconds:      totalDurationSeconds,
		AvgSessionDurationSeconds: avgSessionDurationSeconds,
		ExerciseDistribution:      exerciseDistribution,
		CategoryDistribution:      categoryDistribution,
		PracticeFrequency:         practiceFrequency,
	}, nil
}

// getCategoryDailyPractice retrieves the daily practice data for a specific category
func (h *PracticeSessionHandler) getCategoryDailyPractice(ctx context.Context, categoryId int32, baseWhereClause string, baseParams []interface{}) ([]*pb.PracticeTimePoint, error) {
	// Build the query for category daily practice
	categoryDailyQuery := `
		SELECT 
			date(ps.start_time) as practice_date,
			SUM(strftime('%s', eh.end_time) - strftime('%s', eh.start_time)) as duration
		FROM 
			exercise_history eh
		JOIN 
			practice_sessions ps ON eh.session_id = ps.id
		JOIN 
			exercise_categories ec ON eh.exercise_id = ec.exercise_id
		WHERE 
			ec.category_id = ?
	`

	// Create parameters for query, starting with category ID
	categoryDailyParams := []interface{}{categoryId}

	// Add date range filters if they exist in the base query
	if baseWhereClause != "" {
		// Convert "WHERE ps.start_time >= ?" to "AND ps.start_time >= ?"
		categoryFilter := strings.Replace(baseWhereClause, "WHERE", "AND", 1)
		categoryDailyQuery += categoryFilter
		categoryDailyParams = append(categoryDailyParams, baseParams...)
	}

	// Add grouping and ordering
	categoryDailyQuery += `
		GROUP BY 
			practice_date
		ORDER BY 
			practice_date ASC
	`

	// Execute the query
	categoryDailyRows, err := h.db.QueryContext(ctx, categoryDailyQuery, categoryDailyParams...)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to calculate daily practice for category %d: %v", categoryId, err)
	}
	defer categoryDailyRows.Close()

	var categoryDailyPoints []*pb.PracticeTimePoint

	for categoryDailyRows.Next() {
		var dateStr string
		var durationSeconds int32

		if err := categoryDailyRows.Scan(&dateStr, &durationSeconds); err != nil {
			return nil, status.Errorf(codes.Internal, "failed to parse daily category practice: %v", err)
		}

		// Parse the date string (format: YYYY-MM-DD)
		practiceDate, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to parse date: %v", err)
		}

		point := &pb.PracticeTimePoint{
			Date:            timestamppb.New(practiceDate),
			DurationSeconds: durationSeconds,
		}
		categoryDailyPoints = append(categoryDailyPoints, point)
	}

	if err = categoryDailyRows.Err(); err != nil {
		return nil, status.Errorf(codes.Internal, "error reading daily category practice: %v", err)
	}

	return categoryDailyPoints, nil
}

// Helper method to get exercise details
func (h *PracticeSessionHandler) getExerciseDetails(ctx context.Context, tx *sql.Tx, exerciseId int32) (*pb.Exercise, error) {
	// Get basic exercise info
	var exercise pb.Exercise
	var createdAt, updatedAt time.Time

	err := tx.QueryRowContext(
		ctx,
		"SELECT id, name, description, created_at, updated_at FROM exercises WHERE id = ?",
		exerciseId,
	).Scan(&exercise.Id, &exercise.Name, &exercise.Description, &createdAt, &updatedAt)

	if err == sql.ErrNoRows {
		return nil, status.Errorf(codes.NotFound, "exercise with ID %d not found", exerciseId)
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to retrieve exercise: %v", err)
	}

	exercise.CreatedAt = timestamppb.New(createdAt)
	exercise.UpdatedAt = timestamppb.New(updatedAt)

	// Get associated tag IDs
	tagRows, err := tx.QueryContext(
		ctx,
		"SELECT tag_id FROM exercise_tags WHERE exercise_id = ?",
		exerciseId,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to retrieve exercise tags: %v", err)
	}
	defer tagRows.Close()

	var tagIDs []int32
	for tagRows.Next() {
		var tagID int32
		if err := tagRows.Scan(&tagID); err != nil {
			return nil, status.Errorf(codes.Internal, "failed to parse tag ID: %v", err)
		}
		tagIDs = append(tagIDs, tagID)
	}
	if err = tagRows.Err(); err != nil {
		return nil, status.Errorf(codes.Internal, "error reading exercise tags: %v", err)
	}
	exercise.TagIds = tagIDs

	// Get associated category IDs
	categoryRows, err := tx.QueryContext(
		ctx,
		"SELECT category_id FROM exercise_categories WHERE exercise_id = ?",
		exerciseId,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to retrieve exercise categories: %v", err)
	}
	defer categoryRows.Close()

	var categoryIDs []int32
	for categoryRows.Next() {
		var categoryID int32
		if err := categoryRows.Scan(&categoryID); err != nil {
			return nil, status.Errorf(codes.Internal, "failed to parse category ID: %v", err)
		}
		categoryIDs = append(categoryIDs, categoryID)
	}
	if err = categoryRows.Err(); err != nil {
		return nil, status.Errorf(codes.Internal, "error reading exercise categories: %v", err)
	}
	exercise.CategoryIds = categoryIDs

	return &exercise, nil
}
