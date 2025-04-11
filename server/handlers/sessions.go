package handlers

import (
	"context"
	"database/sql"
	"strconv"
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

	// Return the created session (without exercises for now)
	return &pb.PracticeSession{
		Id:        sessionID,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		Notes:     req.Notes,
		CreatedAt: timestamppb.New(createdAt),
		UpdatedAt: timestamppb.New(updatedAt),
		Exercises: []*pb.SessionExercise{}, // Empty for now
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

	// Query session exercises
	exerciseRows, err := tx.QueryContext(
		ctx,
		`SELECT id, exercise_id, start_time, end_time, bpm, time_signature, notes
         FROM session_exercises
         WHERE session_id = ?
         ORDER BY start_time`,
		req.Id,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to retrieve session exercises: %v", err)
	}
	defer exerciseRows.Close()

	// Parse exercises
	var exercises []*pb.SessionExercise
	exerciseIDMap := make(map[int64]int) // Maps exercise ID to index in exercises slice

	for exerciseRows.Next() {
		var sessionExercise pb.SessionExercise
		var exerciseId int64
		var startTime, endTime time.Time

		err := exerciseRows.Scan(
			&sessionExercise.Id,
			&exerciseId,
			&startTime,
			&endTime,
			&sessionExercise.Bpm,
			&sessionExercise.TimeSignature,
			&sessionExercise.Notes,
		)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to parse session exercise: %v", err)
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
		whereClause += " id IN (SELECT session_id FROM session_exercises WHERE exercise_id = ?)"
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
		session.Exercises = []*pb.SessionExercise{} // Empty for now

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

// AddSessionExercise adds an exercise to a session
func (h *PracticeSessionHandler) AddSessionExercise(ctx context.Context, req *pb.AddSessionExerciseRequest) (*pb.SessionExercise, error) {
	if req.SessionId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "invalid session ID")
	}

	if req.ExerciseId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "invalid exercise ID")
	}

	if req.StartTime == nil || req.EndTime == nil {
		return nil, status.Error(codes.InvalidArgument, "start time and end time are required")
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

	// Check if session exists
	var sessionExists bool
	err = tx.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM practice_sessions WHERE id = ?)", req.SessionId).Scan(&sessionExists)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to check session existence: %v", err)
	}
	if !sessionExists {
		return nil, status.Errorf(codes.NotFound, "session with ID %d not found", req.SessionId)
	}

	// Check if exercise exists
	var exerciseExists bool
	err = tx.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM exercises WHERE id = ?)", req.ExerciseId).Scan(&exerciseExists)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to check exercise existence: %v", err)
	}
	if !exerciseExists {
		return nil, status.Errorf(codes.NotFound, "exercise with ID %d not found", req.ExerciseId)
	}

	// Insert the session exercise
	result, err := tx.ExecContext(
		ctx,
		`INSERT INTO session_exercises (session_id, exercise_id, start_time, end_time, bpm, time_signature, notes) 
         VALUES (?, ?, ?, ?, ?, ?, ?)`,
		req.SessionId, req.ExerciseId, startTime, endTime, req.Bpm, req.TimeSignature, req.Notes,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to add exercise to session: %v", err)
	}

	// Get the session exercise ID
	sessionExerciseId, err := result.LastInsertId()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get session exercise ID: %v", err)
	}

	// Fetch exercise details
	exercise, err := h.getExerciseDetails(ctx, tx, req.ExerciseId)
	if err != nil {
		return nil, err
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to commit transaction: %v", err)
	}

	// Return the created session exercise
	return &pb.SessionExercise{
		Id:            sessionExerciseId,
		SessionId:     req.SessionId,
		ExerciseId:    req.ExerciseId,
		StartTime:     req.StartTime,
		EndTime:       req.EndTime,
		Bpm:           req.Bpm,
		TimeSignature: req.TimeSignature,
		Notes:         req.Notes,
		Exercise:      exercise,
	}, nil
}

// UpdateSessionExercise updates a session exercise
func (h *PracticeSessionHandler) UpdateSessionExercise(ctx context.Context, req *pb.UpdateSessionExerciseRequest) (*pb.SessionExercise, error) {
	if req.Id <= 0 {
		return nil, status.Error(codes.InvalidArgument, "invalid session exercise ID")
	}

	if req.SessionId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "invalid session ID")
	}

	if req.Exercise == nil {
		return nil, status.Error(codes.InvalidArgument, "session exercise data is required")
	}

	// Start a transaction
	tx, err := h.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to start transaction: %v", err)
	}
	defer tx.Rollback() // Rollback if not committed

	// Check if the session exercise exists and belongs to the specified session
	var exists bool
	err = tx.QueryRowContext(
		ctx,
		"SELECT EXISTS(SELECT 1 FROM session_exercises WHERE id = ? AND session_id = ?)",
		req.Id, req.SessionId,
	).Scan(&exists)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to check session exercise existence: %v", err)
	}
	if !exists {
		return nil, status.Errorf(codes.NotFound, "session exercise with ID %d not found in session %d", req.Id, req.SessionId)
	}

	// Parse update mask
	updateStartTime := false
	updateEndTime := false
	updateBpm := false
	updateTimeSignature := false
	updateNotes := false

	if req.UpdateMask == nil || len(req.UpdateMask.Paths) == 0 {
		// If no update mask is provided, update all fields
		updateStartTime = true
		updateEndTime = true
		updateBpm = true
		updateTimeSignature = true
		updateNotes = true
	} else {
		for _, path := range req.UpdateMask.Paths {
			switch path {
			case "start_time":
				updateStartTime = true
			case "end_time":
				updateEndTime = true
			case "bpm":
				updateBpm = true
			case "time_signature":
				updateTimeSignature = true
			case "notes":
				updateNotes = true
			}
		}
	}

	// Validate times if updating
	if updateStartTime && updateEndTime {
		startTime := req.Exercise.StartTime.AsTime()
		endTime := req.Exercise.EndTime.AsTime()
		if startTime.After(endTime) {
			return nil, status.Error(codes.InvalidArgument, "start time cannot be after end time")
		}
	} else if updateStartTime {
		// If only updating start time, check against existing end time
		var endTime time.Time
		err := tx.QueryRowContext(
			ctx,
			"SELECT end_time FROM session_exercises WHERE id = ?",
			req.Id,
		).Scan(&endTime)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to get existing end time: %v", err)
		}

		startTime := req.Exercise.StartTime.AsTime()
		if startTime.After(endTime) {
			return nil, status.Error(codes.InvalidArgument, "start time cannot be after end time")
		}
	} else if updateEndTime {
		// If only updating end time, check against existing start time
		var startTime time.Time
		err := tx.QueryRowContext(
			ctx,
			"SELECT start_time FROM session_exercises WHERE id = ?",
			req.Id,
		).Scan(&startTime)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to get existing start time: %v", err)
		}

		endTime := req.Exercise.EndTime.AsTime()
		if startTime.After(endTime) {
			return nil, status.Error(codes.InvalidArgument, "start time cannot be after end time")
		}
	}

	// Build update SQL
	sql := "UPDATE session_exercises SET"
	params := []interface{}{}
	first := true

	if updateStartTime {
		if req.Exercise.StartTime == nil {
			return nil, status.Error(codes.InvalidArgument, "start time cannot be empty")
		}
		sql += " start_time = ?"
		params = append(params, req.Exercise.StartTime.AsTime())
		first = false
	}

	if updateEndTime {
		if req.Exercise.EndTime == nil {
			return nil, status.Error(codes.InvalidArgument, "end time cannot be empty")
		}
		if !first {
			sql += ","
		}
		sql += " end_time = ?"
		params = append(params, req.Exercise.EndTime.AsTime())
		first = false
	}

	if updateBpm {
		if !first {
			sql += ","
		}
		sql += " bpm = ?"
		params = append(params, req.Exercise.Bpm)
		first = false
	}

	if updateTimeSignature {
		if !first {
			sql += ","
		}
		sql += " time_signature = ?"
		params = append(params, req.Exercise.TimeSignature)
		first = false
	}

	if updateNotes {
		if !first {
			sql += ","
		}
		sql += " notes = ?"
		params = append(params, req.Exercise.Notes)
	}

	if len(params) == 0 {
		return nil, status.Error(codes.InvalidArgument, "no fields to update")
	}

	sql += " WHERE id = ? AND session_id = ?"
	params = append(params, req.Id, req.SessionId)

	// Execute the update
	_, err = tx.ExecContext(ctx, sql, params...)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update session exercise: %v", err)
	}

	// Get the updated session exercise
	var sessionExercise pb.SessionExercise
	var exerciseId int64
	var startTime, endTime time.Time

	err = tx.QueryRowContext(
		ctx,
		`SELECT id, session_id, exercise_id, start_time, end_time, bpm, time_signature, notes
         FROM session_exercises
         WHERE id = ? AND session_id = ?`,
		req.Id, req.SessionId,
	).Scan(
		&sessionExercise.Id,
		&sessionExercise.SessionId,
		&exerciseId,
		&startTime,
		&endTime,
		&sessionExercise.Bpm,
		&sessionExercise.TimeSignature,
		&sessionExercise.Notes,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to retrieve updated session exercise: %v", err)
	}

	sessionExercise.ExerciseId = exerciseId
	sessionExercise.StartTime = timestamppb.New(startTime)
	sessionExercise.EndTime = timestamppb.New(endTime)

	// Fetch exercise details
	exercise, err := h.getExerciseDetails(ctx, tx, exerciseId)
	if err != nil {
		return nil, err
	}
	sessionExercise.Exercise = exercise

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to commit transaction: %v", err)
	}

	return &sessionExercise, nil
}

// DeleteSessionExercise deletes a session exercise
func (h *PracticeSessionHandler) DeleteSessionExercise(ctx context.Context, req *pb.DeleteSessionExerciseRequest) (*emptypb.Empty, error) {
	if req.Id <= 0 {
		return nil, status.Error(codes.InvalidArgument, "invalid session exercise ID")
	}

	// Check if the session exercise exists
	var exists bool
	err := h.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM session_exercises WHERE id = ?)", req.Id).Scan(&exists)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to check session exercise existence: %v", err)
	}
	if !exists {
		return nil, status.Errorf(codes.NotFound, "session exercise with ID %d not found", req.Id)
	}

	// Delete the session exercise
	_, err = h.db.ExecContext(ctx, "DELETE FROM session_exercises WHERE id = ?", req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete session exercise: %v", err)
	}

	return &emptypb.Empty{}, nil
}

// Helper method to get exercise details
func (h *PracticeSessionHandler) getExerciseDetails(ctx context.Context, tx *sql.Tx, exerciseId int64) (*pb.Exercise, error) {
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

	var tagIDs []int64
	for tagRows.Next() {
		var tagID int64
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

	var categoryIDs []int64
	for categoryRows.Next() {
		var categoryID int64
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
