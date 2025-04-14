package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	pb "github.com/Zach-Johnson/tempus/proto/api/v1/tempus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ExerciseHistoryHandler implements the ExerciseHistoryService gRPC service
type ExerciseHistoryHandler struct {
	pb.UnimplementedExerciseHistoryServiceServer
	db *sql.DB
}

// NewExerciseHistoryHandler creates a new ExerciseHistoryHandler
func NewExerciseHistoryHandler(db *sql.DB) *ExerciseHistoryHandler {
	return &ExerciseHistoryHandler{db: db}
}

// CreateExerciseHistory creates a new exercise history entry
func (h *ExerciseHistoryHandler) CreateExerciseHistory(ctx context.Context, req *pb.CreateExerciseHistoryRequest) (*pb.ExerciseHistory, error) {
	if req.SessionId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "invalid session ID")
	}

	if req.ExerciseId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "exercise ID is required")
	}

	if req.StartTime == nil || req.EndTime == nil {
		return nil, status.Error(codes.InvalidArgument, "start time and end time are required")
	}

	startTime := req.StartTime.AsTime()
	endTime := req.EndTime.AsTime()

	if startTime.After(endTime) {
		return nil, status.Error(codes.InvalidArgument, "start time cannot be after end time")
	}

	// Validate rating if provided
	if req.Rating < 0 || req.Rating > 5 {
		return nil, status.Error(codes.InvalidArgument, "rating must be between 0 and 5")
	}

	// Start a transaction
	tx, err := h.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to start transaction: %v", err)
	}
	defer tx.Rollback() // Rollback if not committed

	// Check if exercise exists
	var exerciseExists bool
	err = tx.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM exercises WHERE id = ?)", req.ExerciseId).Scan(&exerciseExists)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to check exercise existence: %v", err)
	}
	if !exerciseExists {
		return nil, status.Errorf(codes.NotFound, "exercise with ID %d not found", req.ExerciseId)
	}

	// Check if session exists
	var sessionExists bool
	err = tx.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM practice_sessions WHERE id = ?)", req.SessionId).Scan(&sessionExists)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to check session existence: %v", err)
	}
	if !sessionExists {
		return nil, status.Errorf(codes.NotFound, "session with ID %d not found", req.SessionId)
	}

	bpmJSON, err := json.Marshal(req.Bpms)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to marshal BPM values: %v", err)
	}

	// Insert the exercise history entry
	result, err := tx.ExecContext(
		ctx,
		`INSERT INTO exercise_history (exercise_id, session_id, start_time, end_time, bpms, time_signature, notes, rating) 
         VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		req.ExerciseId, req.SessionId, startTime, endTime, bpmJSON, req.TimeSignature, req.Notes, req.Rating,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create exercise history entry: %v", err)
	}

	// Get the history entry ID
	historyId, err := result.LastInsertId()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get history entry ID: %v", err)
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

	// Return the created history entry
	return &pb.ExerciseHistory{
		Id:            int32(historyId),
		ExerciseId:    req.ExerciseId,
		StartTime:     req.StartTime,
		EndTime:       req.EndTime,
		Bpms:          req.Bpms,
		TimeSignature: req.TimeSignature,
		Notes:         req.Notes,
		Rating:        req.Rating,
		Exercise:      exercise,
	}, nil
}

// GetExerciseHistory retrieves an exercise history entry by ID
func (h *ExerciseHistoryHandler) GetExerciseHistory(ctx context.Context, req *pb.GetExerciseHistoryRequest) (*pb.ExerciseHistory, error) {
	if req.Id <= 0 {
		return nil, status.Error(codes.InvalidArgument, "invalid history entry ID")
	}

	// Start a transaction to ensure consistency
	tx, err := h.db.BeginTx(ctx, &sql.TxOptions{ReadOnly: true})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to start transaction: %v", err)
	}
	defer tx.Rollback() // Rollback if not committed

	// Query the history entry
	var history pb.ExerciseHistory
	var exerciseId int32
	var startTime, endTime time.Time
	var bpmJSON string

	err = tx.QueryRowContext(
		ctx,
		`SELECT id, exercise_id, session_id, start_time, end_time, bpms, time_signature, notes, rating
     FROM exercise_history
     WHERE id = ?`,
		req.Id,
	).Scan(
		&history.Id,
		&exerciseId,
		&history.SessionId,
		&startTime,
		&endTime,
		&bpmJSON,
		&history.TimeSignature,
		&history.Notes,
		&history.Rating,
	)

	if bpmJSON != "" {
		var bpms []int32
		if err := json.Unmarshal([]byte(bpmJSON), &bpms); err != nil {
			return nil, status.Errorf(codes.Internal, "failed to unmarshal BPM values: %v", err)
		}
		history.Bpms = bpms
	}

	if err == sql.ErrNoRows {
		return nil, status.Errorf(codes.NotFound, "exercise history entry with ID %d not found", req.Id)
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to retrieve exercise history: %v", err)
	}

	history.ExerciseId = exerciseId
	history.StartTime = timestamppb.New(startTime)
	history.EndTime = timestamppb.New(endTime)

	// Fetch exercise details
	exercise, err := h.getExerciseDetails(ctx, tx, exerciseId)
	if err != nil {
		return nil, err
	}
	history.Exercise = exercise

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to commit transaction: %v", err)
	}

	return &history, nil
}

// ListExerciseHistory lists exercise history entries with optional filtering and pagination
func (h *ExerciseHistoryHandler) ListExerciseHistory(ctx context.Context, req *pb.ListExerciseHistoryRequest) (*pb.ListExerciseHistoryResponse, error) {
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
        SELECT id, exercise_id, session_id, start_time, end_time, bpms, time_signature, notes, rating
        FROM exercise_history
    `
	countQuery := `
        SELECT COUNT(*)
        FROM exercise_history
    `

	var whereClause string
	var queryParams []interface{}

	// Add filter by exercise if provided
	if req.ExerciseId > 0 {
		whereClause = " WHERE exercise_id = ?"
		queryParams = append(queryParams, req.ExerciseId)
	}

	// Add filter by date range if provided
	if req.StartDate != nil {
		if whereClause == "" {
			whereClause = " WHERE"
		} else {
			whereClause += " AND"
		}
		whereClause += " start_time >= ?"
		queryParams = append(queryParams, req.StartDate.AsTime())
	}

	if req.EndDate != nil {
		if whereClause == "" {
			whereClause = " WHERE"
		} else {
			whereClause += " AND"
		}
		whereClause += " end_time <= ?"
		queryParams = append(queryParams, req.EndDate.AsTime())
	}

	// Add filter by session if provided
	if req.SessionId > 0 {
		if whereClause == "" {
			whereClause = " WHERE"
		} else {
			whereClause += " AND"
		}
		whereClause += " session_id = ?"
		queryParams = append(queryParams, req.SessionId)
	}

	// Add order by, limit, and offset
	fullQuery := baseQuery + whereClause + " ORDER BY start_time DESC LIMIT ? OFFSET ?"
	queryParams = append(queryParams, pageSize+1, offset) // Query one more to check if there are more pages

	// Start a transaction to ensure consistency
	tx, err := h.db.BeginTx(ctx, &sql.TxOptions{ReadOnly: true})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to start transaction: %v", err)
	}
	defer tx.Rollback() // Rollback if not committed

	// Query total count
	var totalCount int32
	countQueryParams := make([]interface{}, len(queryParams)-2) // Exclude limit and offset
	copy(countQueryParams, queryParams[:len(queryParams)-2])

	err = tx.QueryRowContext(ctx, countQuery+whereClause, countQueryParams...).Scan(&totalCount)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to count exercise history entries: %v", err)
	}

	// Query history entries
	rows, err := tx.QueryContext(ctx, fullQuery, queryParams...)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list exercise history: %v", err)
	}
	defer rows.Close()

	// Parse the results
	historyEntries := make([]*pb.ExerciseHistory, 0, pageSize)
	count := 0
	hasMorePages := false
	exerciseIDs := make(map[int32]bool)

	for rows.Next() {
		if count >= pageSize {
			hasMorePages = true
			break
		}

		var history pb.ExerciseHistory
		var exerciseId int32
		var sessionID int32
		var startTime, endTime time.Time
		var bpmJSON string

		err := rows.Scan(
			&history.Id,
			&exerciseId,
			&sessionID,
			&startTime,
			&endTime,
			&bpmJSON,
			&history.TimeSignature,
			&history.Notes,
			&history.Rating,
		)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to parse exercise history: %v", err)
		}

		if bpmJSON != "" {
			var bpms []int32
			if err := json.Unmarshal([]byte(bpmJSON), &bpms); err != nil {
				return nil, status.Errorf(codes.Internal, "failed to unmarshal BPM values: %v", err)
			}
			history.Bpms = bpms
		}

		history.ExerciseId = exerciseId
		history.SessionId = sessionID
		history.StartTime = timestamppb.New(startTime)
		history.EndTime = timestamppb.New(endTime)

		exerciseIDs[exerciseId] = true
		historyEntries = append(historyEntries, &history)
		count++
	}

	if err = rows.Err(); err != nil {
		return nil, status.Errorf(codes.Internal, "error reading exercise history: %v", err)
	}

	// Fetch exercise details
	if len(historyEntries) > 0 {
		exercises, err := h.fetchExercises(ctx, tx, exerciseIDs)
		if err != nil {
			return nil, err
		}

		// Associate exercises with history entries
		for _, entry := range historyEntries {
			if exercise, ok := exercises[entry.ExerciseId]; ok {
				entry.Exercise = exercise
			}
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to commit transaction: %v", err)
	}

	// Calculate next page token
	nextPageToken := ""
	if hasMorePages {
		nextPageToken = strconv.Itoa(offset + pageSize)
	}

	return &pb.ListExerciseHistoryResponse{
		HistoryEntries: historyEntries,
		NextPageToken:  nextPageToken,
		TotalCount:     totalCount,
	}, nil
}

// UpdateExerciseHistory updates an exercise history entry
func (h *ExerciseHistoryHandler) UpdateExerciseHistory(ctx context.Context, req *pb.UpdateExerciseHistoryRequest) (*pb.ExerciseHistory, error) {
	if req.Id <= 0 {
		return nil, status.Error(codes.InvalidArgument, "invalid history entry ID")
	}

	if req.History == nil {
		return nil, status.Error(codes.InvalidArgument, "history data is required")
	}

	// Check if the history entry exists
	var exists bool
	err := h.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM exercise_history WHERE id = ?)", req.Id).Scan(&exists)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to check history entry existence: %v", err)
	}
	if !exists {
		return nil, status.Errorf(codes.NotFound, "history entry with ID %d not found", req.Id)
	}

	// Validate rating if provided
	if req.History.Rating < 0 || req.History.Rating > 5 {
		return nil, status.Error(codes.InvalidArgument, "rating must be between 0 and 5")
	}

	// Parse update mask
	updateStartTime := false
	updateEndTime := false
	updateBpm := false
	updateTimeSignature := false
	updateNotes := false
	updateRating := false

	if req.UpdateMask == nil || len(req.UpdateMask.Paths) == 0 {
		// If no update mask is provided, update all fields
		updateStartTime = true
		updateEndTime = true
		updateBpm = true
		updateTimeSignature = true
		updateNotes = true
		updateRating = true
	} else {
		for _, path := range req.UpdateMask.Paths {
			switch path {
			case "start_time":
				updateStartTime = true
			case "end_time":
				updateEndTime = true
			case "bpms":
				updateBpm = true
			case "time_signature":
				updateTimeSignature = true
			case "notes":
				updateNotes = true
			case "rating":
				updateRating = true
			}
		}
	}

	// Validate times if updating
	if updateStartTime && updateEndTime {
		startTime := req.History.StartTime.AsTime()
		endTime := req.History.EndTime.AsTime()
		if startTime.After(endTime) {
			return nil, status.Error(codes.InvalidArgument, "start time cannot be after end time")
		}
	} else if updateStartTime {
		// If only updating start time, check against existing end time
		var endTime time.Time
		err := h.db.QueryRowContext(ctx, "SELECT end_time FROM exercise_history WHERE id = ?", req.Id).Scan(&endTime)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to get existing end time: %v", err)
		}

		startTime := req.History.StartTime.AsTime()
		if startTime.After(endTime) {
			return nil, status.Error(codes.InvalidArgument, "start time cannot be after end time")
		}
	} else if updateEndTime {
		// If only updating end time, check against existing start time
		var startTime time.Time
		err := h.db.QueryRowContext(ctx, "SELECT start_time FROM exercise_history WHERE id = ?", req.Id).Scan(&startTime)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to get existing start time: %v", err)
		}

		endTime := req.History.EndTime.AsTime()
		if startTime.After(endTime) {
			return nil, status.Error(codes.InvalidArgument, "start time cannot be after end time")
		}
	}

	// Build update SQL
	sql := "UPDATE exercise_history SET"
	params := []interface{}{}
	first := true

	if updateStartTime {
		if req.History.StartTime == nil {
			return nil, status.Error(codes.InvalidArgument, "start time cannot be empty")
		}
		sql += " start_time = ?"
		params = append(params, req.History.StartTime.AsTime())
		first = false
	}

	if updateEndTime {
		if req.History.EndTime == nil {
			return nil, status.Error(codes.InvalidArgument, "end time cannot be empty")
		}
		if !first {
			sql += ","
		}
		sql += " end_time = ?"
		params = append(params, req.History.EndTime.AsTime())
		first = false
	}

	if updateBpm {
		bpmJSON, err := json.Marshal(req.History.Bpms)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to marshal BPM values: %v", err)
		}

		if !first {
			sql += ","
		}
		sql += " bpms = ?"
		params = append(params, bpmJSON)
		first = false
	}

	if updateTimeSignature {
		if !first {
			sql += ","
		}
		sql += " time_signature = ?"
		params = append(params, req.History.TimeSignature)
		first = false
	}

	if updateNotes {
		if !first {
			sql += ","
		}
		sql += " notes = ?"
		params = append(params, req.History.Notes)
		first = false
	}

	if updateRating {
		if !first {
			sql += ","
		}
		sql += " rating = ?"
		params = append(params, req.History.Rating)
	}

	tx, err := h.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to start transaction: %v", err)
	}
	defer tx.Rollback() // Rollback if not committed

	if len(params) == 0 {
		return nil, status.Error(codes.InvalidArgument, "no fields to update")
	}

	sql += " WHERE id = ?"
	params = append(params, req.Id)

	// Execute the update
	_, err = tx.ExecContext(ctx, sql, params...)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update exercise history: %v", err)
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to commit transaction: %v", err)
	}

	// Retrieve the updated history entry
	return h.GetExerciseHistory(ctx, &pb.GetExerciseHistoryRequest{Id: req.Id})
}

// DeleteExerciseHistory deletes an exercise history entry
func (h *ExerciseHistoryHandler) DeleteExerciseHistory(ctx context.Context, req *pb.DeleteExerciseHistoryRequest) (*emptypb.Empty, error) {
	if req.Id <= 0 {
		return nil, status.Error(codes.InvalidArgument, "invalid history entry ID")
	}

	// Check if the history entry exists
	var exists bool
	err := h.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM exercise_history WHERE id = ?)", req.Id).Scan(&exists)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to check history entry existence: %v", err)
	}
	if !exists {
		return nil, status.Errorf(codes.NotFound, "history entry with ID %d not found", req.Id)
	}

	// Delete the history entry
	_, err = h.db.ExecContext(ctx, "DELETE FROM exercise_history WHERE id = ?", req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete exercise history entry: %v", err)
	}

	return &emptypb.Empty{}, nil
}

// Helper method to get exercise details
func (h *ExerciseHistoryHandler) getExerciseDetails(ctx context.Context, tx *sql.Tx, exerciseId int32) (*pb.Exercise, error) {
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

// Helper method to fetch multiple exercises
func (h *ExerciseHistoryHandler) fetchExercises(ctx context.Context, tx *sql.Tx, exerciseIDs map[int32]bool) (map[int32]*pb.Exercise, error) {
	// Convert map keys to slice
	ids := make([]int32, 0, len(exerciseIDs))
	for id := range exerciseIDs {
		ids = append(ids, id)
	}

	// Build placeholders for the IN clause
	placeholders := make([]string, len(ids))
	args := make([]interface{}, len(ids))
	for i, id := range ids {
		placeholders[i] = "?"
		args[i] = id
	}

	// Query basic exercise info
	query := fmt.Sprintf(
		`SELECT id, name, description, created_at, updated_at
         FROM exercises
         WHERE id IN (%s)`,
		strings.Join(placeholders, ","),
	)

	rows, err := tx.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to fetch exercises: %v", err)
	}
	defer rows.Close()

	// Parse exercises
	result := make(map[int32]*pb.Exercise)
	for rows.Next() {
		var exercise pb.Exercise
		var createdAt, updatedAt time.Time

		err := rows.Scan(
			&exercise.Id,
			&exercise.Name,
			&exercise.Description,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to parse exercise: %v", err)
		}

		exercise.CreatedAt = timestamppb.New(createdAt)
		exercise.UpdatedAt = timestamppb.New(updatedAt)
		result[exercise.Id] = &exercise
	}

	if err = rows.Err(); err != nil {
		return nil, status.Errorf(codes.Internal, "error reading exercises: %v", err)
	}

	// Fetch tags, categories for these exercises
	if len(result) > 0 {
		if err := h.fetchExerciseRelations(ctx, tx, result, ids); err != nil {
			return nil, err
		}
	}

	return result, nil
}

// Helper to fetch tags and categories for exercises
func (h *ExerciseHistoryHandler) fetchExerciseRelations(ctx context.Context, tx *sql.Tx, exercises map[int32]*pb.Exercise, exerciseIDs []int32) error {
	// Build placeholders for the IN clause
	placeholders := make([]string, len(exerciseIDs))
	args := make([]interface{}, len(exerciseIDs))
	for i, id := range exerciseIDs {
		placeholders[i] = "?"
		args[i] = id
	}

	// Fetch tags
	tagQuery := fmt.Sprintf(
		"SELECT exercise_id, tag_id FROM exercise_tags WHERE exercise_id IN (%s)",
		strings.Join(placeholders, ","),
	)
	tagRows, err := tx.QueryContext(ctx, tagQuery, args...)
	if err != nil {
		return status.Errorf(codes.Internal, "failed to fetch exercise tags: %v", err)
	}
	defer tagRows.Close()

	for tagRows.Next() {
		var exerciseID, tagID int32
		if err := tagRows.Scan(&exerciseID, &tagID); err != nil {
			return status.Errorf(codes.Internal, "failed to parse exercise tag: %v", err)
		}
		if exercise, ok := exercises[exerciseID]; ok {
			exercise.TagIds = append(exercise.TagIds, tagID)
		}
	}
	if err = tagRows.Err(); err != nil {
		return status.Errorf(codes.Internal, "error reading exercise tags: %v", err)
	}

	// Fetch categories
	catQuery := fmt.Sprintf(
		"SELECT exercise_id, category_id FROM exercise_categories WHERE exercise_id IN (%s)",
		strings.Join(placeholders, ","),
	)
	catRows, err := tx.QueryContext(ctx, catQuery, args...)
	if err != nil {
		return status.Errorf(codes.Internal, "failed to fetch exercise categories: %v", err)
	}
	defer catRows.Close()

	for catRows.Next() {
		var exerciseID, categoryID int32
		if err := catRows.Scan(&exerciseID, &categoryID); err != nil {
			return status.Errorf(codes.Internal, "failed to parse exercise category: %v", err)
		}
		if exercise, ok := exercises[exerciseID]; ok {
			exercise.CategoryIds = append(exercise.CategoryIds, categoryID)
		}
	}
	if err = catRows.Err(); err != nil {
		return status.Errorf(codes.Internal, "error reading exercise categories: %v", err)
	}

	return nil
}
