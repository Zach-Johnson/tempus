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

// ExerciseHandler implements the ExerciseService gRPC service
type ExerciseHandler struct {
	pb.UnimplementedExerciseServiceServer
	db *sql.DB
}

// NewExerciseHandler creates a new ExerciseHandler
func NewExerciseHandler(db *sql.DB) *ExerciseHandler {
	return &ExerciseHandler{db: db}
}

// CreateExercise creates a new exercise
func (h *ExerciseHandler) CreateExercise(ctx context.Context, req *pb.CreateExerciseRequest) (*pb.Exercise, error) {
	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "exercise name is required")
	}

	// Start a transaction
	tx, err := h.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to start transaction: %v", err)
	}
	defer tx.Rollback() // Rollback if not committed

	// Insert the exercise
	result, err := tx.ExecContext(
		ctx,
		"INSERT INTO exercises (name, description) VALUES (?, ?)",
		req.Name, req.Description,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create exercise: %v", err)
	}

	// Get the exercise ID
	id, err := result.LastInsertId()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get exercise ID: %v", err)
	}

	// Associate exercise with tags if provided
	for _, tagID := range req.TagIds {
		// Check if tag exists
		var exists bool
		err := tx.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM tags WHERE id = ?)", tagID).Scan(&exists)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to check tag existence: %v", err)
		}
		if !exists {
			return nil, status.Errorf(codes.NotFound, "tag with ID %d not found", tagID)
		}

		// Create exercise-tag relationship
		_, err = tx.ExecContext(
			ctx,
			"INSERT INTO exercise_tags (exercise_id, tag_id) VALUES (?, ?)",
			id, tagID,
		)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to associate exercise with tag: %v", err)
		}
	}

	// Associate exercise with categories if provided
	for _, categoryID := range req.CategoryIds {
		// Check if category exists
		var exists bool
		err := tx.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM categories WHERE id = ?)", categoryID).Scan(&exists)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to check category existence: %v", err)
		}
		if !exists {
			return nil, status.Errorf(codes.NotFound, "category with ID %d not found", categoryID)
		}

		// Create exercise-category relationship
		_, err = tx.ExecContext(
			ctx,
			"INSERT INTO exercise_categories (exercise_id, category_id) VALUES (?, ?)",
			id, categoryID,
		)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to associate exercise with category: %v", err)
		}
	}

	// Add images if provided
	images := make([]*pb.ExerciseImage, 0, len(req.Images))
	for _, imageReq := range req.Images {
		result, err := tx.ExecContext(
			ctx,
			`INSERT INTO exercise_images (exercise_id, image_data, filename, mime_type, description) 
             VALUES (?, ?, ?, ?, ?)`,
			id, imageReq.ImageData, imageReq.Filename, imageReq.MimeType, imageReq.Description,
		)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to add image to exercise: %v", err)
		}

		imageID, err := result.LastInsertId()
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to get image ID: %v", err)
		}

		var imageCreatedAt time.Time
		err = tx.QueryRowContext(
			ctx,
			"SELECT created_at FROM exercise_images WHERE id = ?",
			imageID,
		).Scan(&imageCreatedAt)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to get image creation time: %v", err)
		}

		images = append(images, &pb.ExerciseImage{
			Id:          int32(imageID),
			ExerciseId:  int32(id),
			ImageData:   imageReq.ImageData,
			Filename:    imageReq.Filename,
			MimeType:    imageReq.MimeType,
			Description: imageReq.Description,
			CreatedAt:   timestamppb.New(imageCreatedAt),
		})
	}

	// Add links if provided
	links := make([]*pb.ExerciseLink, 0, len(req.Links))
	for _, linkReq := range req.Links {
		result, err := tx.ExecContext(
			ctx,
			`INSERT INTO exercise_links (exercise_id, url, description) 
             VALUES (?, ?, ?)`,
			id, linkReq.Url, linkReq.Description,
		)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to add link to exercise: %v", err)
		}

		linkID, err := result.LastInsertId()
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to get link ID: %v", err)
		}

		var linkCreatedAt time.Time
		err = tx.QueryRowContext(
			ctx,
			"SELECT created_at FROM exercise_links WHERE id = ?",
			linkID,
		).Scan(&linkCreatedAt)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to get link creation time: %v", err)
		}

		links = append(links, &pb.ExerciseLink{
			Id:          int32(linkID),
			ExerciseId:  int32(id),
			Url:         linkReq.Url,
			Description: linkReq.Description,
			CreatedAt:   timestamppb.New(linkCreatedAt),
		})
	}

	// Get the exercise creation and update times
	var createdAt, updatedAt time.Time
	err = tx.QueryRowContext(
		ctx,
		"SELECT created_at, updated_at FROM exercises WHERE id = ?",
		id,
	).Scan(&createdAt, &updatedAt)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get exercise creation time: %v", err)
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to commit transaction: %v", err)
	}

	// Return the created exercise
	return &pb.Exercise{
		Id:          int32(id),
		Name:        req.Name,
		Description: req.Description,
		CreatedAt:   timestamppb.New(createdAt),
		UpdatedAt:   timestamppb.New(updatedAt),
		TagIds:      req.TagIds,
		CategoryIds: req.CategoryIds,
		Images:      images,
		Links:       links,
	}, nil
}

// GetExercise retrieves an exercise by ID
func (h *ExerciseHandler) GetExercise(ctx context.Context, req *pb.GetExerciseRequest) (*pb.Exercise, error) {
	if req.Id <= 0 {
		return nil, status.Error(codes.InvalidArgument, "invalid exercise ID")
	}

	// Query the exercise
	var exercise pb.Exercise
	var createdAt, updatedAt time.Time

	err := h.db.QueryRowContext(
		ctx,
		"SELECT id, name, description, created_at, updated_at FROM exercises WHERE id = ?",
		req.Id,
	).Scan(&exercise.Id, &exercise.Name, &exercise.Description, &createdAt, &updatedAt)

	if err == sql.ErrNoRows {
		return nil, status.Errorf(codes.NotFound, "exercise with ID %d not found", req.Id)
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to retrieve exercise: %v", err)
	}

	exercise.CreatedAt = timestamppb.New(createdAt)
	exercise.UpdatedAt = timestamppb.New(updatedAt)

	// Get associated tag IDs
	tagRows, err := h.db.QueryContext(
		ctx,
		"SELECT tag_id FROM exercise_tags WHERE exercise_id = ?",
		req.Id,
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
	categoryRows, err := h.db.QueryContext(
		ctx,
		"SELECT category_id FROM exercise_categories WHERE exercise_id = ?",
		req.Id,
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

	// Get images
	imageRows, err := h.db.QueryContext(
		ctx,
		`SELECT id, image_data, filename, mime_type, description, created_at 
         FROM exercise_images WHERE exercise_id = ?`,
		req.Id,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to retrieve exercise images: %v", err)
	}
	defer imageRows.Close()

	var images []*pb.ExerciseImage
	for imageRows.Next() {
		var image pb.ExerciseImage
		var imageCreatedAt time.Time
		if err := imageRows.Scan(
			&image.Id, &image.ImageData, &image.Filename, &image.MimeType, &image.Description, &imageCreatedAt,
		); err != nil {
			return nil, status.Errorf(codes.Internal, "failed to parse exercise image: %v", err)
		}
		image.ExerciseId = req.Id
		image.CreatedAt = timestamppb.New(imageCreatedAt)
		images = append(images, &image)
	}
	if err = imageRows.Err(); err != nil {
		return nil, status.Errorf(codes.Internal, "error reading exercise images: %v", err)
	}
	exercise.Images = images

	// Get links
	linkRows, err := h.db.QueryContext(
		ctx,
		`SELECT id, url, description, created_at 
         FROM exercise_links WHERE exercise_id = ?`,
		req.Id,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to retrieve exercise links: %v", err)
	}
	defer linkRows.Close()

	var links []*pb.ExerciseLink
	for linkRows.Next() {
		var link pb.ExerciseLink
		var linkCreatedAt time.Time
		if err := linkRows.Scan(&link.Id, &link.Url, &link.Description, &linkCreatedAt); err != nil {
			return nil, status.Errorf(codes.Internal, "failed to parse exercise link: %v", err)
		}
		link.ExerciseId = req.Id
		link.CreatedAt = timestamppb.New(linkCreatedAt)
		links = append(links, &link)
	}
	if err = linkRows.Err(); err != nil {
		return nil, status.Errorf(codes.Internal, "error reading exercise links: %v", err)
	}
	exercise.Links = links

	return &exercise, nil
}

// ListExercises lists exercises with optional filtering and pagination
func (h *ExerciseHandler) ListExercises(ctx context.Context, req *pb.ListExercisesRequest) (*pb.ListExercisesResponse, error) {
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
        SELECT DISTINCT e.id, e.name, e.description, e.created_at, e.updated_at
        FROM exercises e
    `
	countQuery := `
        SELECT COUNT(DISTINCT e.id)
        FROM exercises e
    `

	var whereClause string
	var queryParams []interface{}

	if req.CategoryId > 0 && req.TagId > 0 {
		// Filter by both category and tag
		whereClause = `
            WHERE e.id IN (
                SELECT ec.exercise_id 
                FROM exercise_categories ec 
                WHERE ec.category_id = ?
            ) AND e.id IN (
                SELECT et.exercise_id 
                FROM exercise_tags et 
                WHERE et.tag_id = ?
            )
        `
		queryParams = append(queryParams, req.CategoryId, req.TagId)
	} else if req.CategoryId > 0 {
		// Filter by category only
		whereClause = `
            WHERE e.id IN (
                SELECT ec.exercise_id 
                FROM exercise_categories ec 
                WHERE ec.category_id = ?
            )
        `
		queryParams = append(queryParams, req.CategoryId)
	} else if req.TagId > 0 {
		// Filter by tag only
		whereClause = `
            WHERE e.id IN (
                SELECT et.exercise_id 
                FROM exercise_tags et 
                WHERE et.tag_id = ?
            )
        `
		queryParams = append(queryParams, req.TagId)
	}

	// Add order by, limit, and offset
	fullQuery := baseQuery + whereClause + " ORDER BY e.name LIMIT ? OFFSET ?"
	queryParams = append(queryParams, pageSize+1, offset) // Query one more to check if there are more pages

	// Query total count
	var totalCount int32
	countQueryParams := make([]interface{}, len(queryParams)-2) // Exclude limit and offset
	copy(countQueryParams, queryParams[:len(queryParams)-2])

	err := h.db.QueryRowContext(ctx, countQuery+whereClause, countQueryParams...).Scan(&totalCount)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to count exercises: %v", err)
	}

	// Query exercises
	rows, err := h.db.QueryContext(ctx, fullQuery, queryParams...)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list exercises: %v", err)
	}
	defer rows.Close()

	// Parse the results
	exercises := make([]*pb.Exercise, 0, pageSize)
	count := 0
	hasMorePages := false
	var exerciseIDs []int32

	for rows.Next() {
		if count >= pageSize {
			hasMorePages = true
			break
		}

		var exercise pb.Exercise
		var createdAt, updatedAt time.Time

		err := rows.Scan(&exercise.Id, &exercise.Name, &exercise.Description, &createdAt, &updatedAt)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to parse exercise: %v", err)
		}

		exercise.CreatedAt = timestamppb.New(createdAt)
		exercise.UpdatedAt = timestamppb.New(updatedAt)
		exercises = append(exercises, &exercise)
		exerciseIDs = append(exerciseIDs, exercise.Id)
		count++
	}

	if err = rows.Err(); err != nil {
		return nil, status.Errorf(codes.Internal, "error reading exercises: %v", err)
	}

	// For each exercise, get tags, categories, images, and links
	if len(exerciseIDs) > 0 {
		if err := h.addRelatedData(ctx, exercises); err != nil {
			return nil, err
		}
	}

	// Calculate next page token
	nextPageToken := ""
	if hasMorePages {
		nextPageToken = strconv.Itoa(offset + pageSize)
	}

	return &pb.ListExercisesResponse{
		Exercises:     exercises,
		NextPageToken: nextPageToken,
		TotalCount:    totalCount,
	}, nil
}

// addRelatedData adds tags, categories, images, and links to the exercises
func (h *ExerciseHandler) addRelatedData(ctx context.Context, exercises []*pb.Exercise) error {
	// Map for quick lookup of exercises by ID
	exerciseMap := make(map[int32]*pb.Exercise)
	exerciseIDs := make([]interface{}, 0, len(exercises))

	// Build query params and map
	placeholders := ""
	for i, ex := range exercises {
		exerciseMap[ex.Id] = ex
		exerciseIDs = append(exerciseIDs, ex.Id)

		if i > 0 {
			placeholders += ","
		}
		placeholders += "?"
	}

	// Get tags for all exercises
	tagQuery := "SELECT exercise_id, tag_id FROM exercise_tags WHERE exercise_id IN (" + placeholders + ")"
	tagRows, err := h.db.QueryContext(ctx, tagQuery, exerciseIDs...)
	if err != nil {
		return status.Errorf(codes.Internal, "failed to retrieve exercise tags: %v", err)
	}
	defer tagRows.Close()

	for tagRows.Next() {
		var exerciseID, tagID int32
		if err := tagRows.Scan(&exerciseID, &tagID); err != nil {
			return status.Errorf(codes.Internal, "failed to parse exercise tag: %v", err)
		}

		if exercise, ok := exerciseMap[exerciseID]; ok {
			exercise.TagIds = append(exercise.TagIds, tagID)
		}
	}
	if err = tagRows.Err(); err != nil {
		return status.Errorf(codes.Internal, "error reading exercise tags: %v", err)
	}

	// Get categories for all exercises
	catQuery := "SELECT exercise_id, category_id FROM exercise_categories WHERE exercise_id IN (" + placeholders + ")"
	catRows, err := h.db.QueryContext(ctx, catQuery, exerciseIDs...)
	if err != nil {
		return status.Errorf(codes.Internal, "failed to retrieve exercise categories: %v", err)
	}
	defer catRows.Close()

	for catRows.Next() {
		var exerciseID, categoryID int32
		if err := catRows.Scan(&exerciseID, &categoryID); err != nil {
			return status.Errorf(codes.Internal, "failed to parse exercise category: %v", err)
		}

		if exercise, ok := exerciseMap[exerciseID]; ok {
			exercise.CategoryIds = append(exercise.CategoryIds, categoryID)
		}
	}
	if err = catRows.Err(); err != nil {
		return status.Errorf(codes.Internal, "error reading exercise categories: %v", err)
	}

	// Get images for all exercises (optional, as these can be large)
	// We could skip this for list calls and only include them in GetExercise

	// Get links for all exercises
	linkQuery := "SELECT id, exercise_id, url, description, created_at FROM exercise_links WHERE exercise_id IN (" + placeholders + ")"
	linkRows, err := h.db.QueryContext(ctx, linkQuery, exerciseIDs...)
	if err != nil {
		return status.Errorf(codes.Internal, "failed to retrieve exercise links: %v", err)
	}
	defer linkRows.Close()

	for linkRows.Next() {
		var link pb.ExerciseLink
		var linkCreatedAt time.Time

		if err := linkRows.Scan(&link.Id, &link.ExerciseId, &link.Url, &link.Description, &linkCreatedAt); err != nil {
			return status.Errorf(codes.Internal, "failed to parse exercise link: %v", err)
		}

		link.CreatedAt = timestamppb.New(linkCreatedAt)

		if exercise, ok := exerciseMap[link.ExerciseId]; ok {
			exercise.Links = append(exercise.Links, &link)
		}
	}
	if err = linkRows.Err(); err != nil {
		return status.Errorf(codes.Internal, "error reading exercise links: %v", err)
	}

	return nil
}

// UpdateExercise updates an exercise
func (h *ExerciseHandler) UpdateExercise(ctx context.Context, req *pb.UpdateExerciseRequest) (*pb.Exercise, error) {
	if req.Id <= 0 {
		return nil, status.Error(codes.InvalidArgument, "invalid exercise ID")
	}

	if req.Exercise == nil {
		return nil, status.Error(codes.InvalidArgument, "exercise data is required")
	}

	// Start transaction
	tx, err := h.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to start transaction: %v", err)
	}
	defer tx.Rollback()

	// Check if exercise exists
	var exists bool
	err = tx.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM exercises WHERE id = ?)", req.Id).Scan(&exists)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to check exercise existence: %v", err)
	}
	if !exists {
		return nil, status.Errorf(codes.NotFound, "exercise with ID %d not found", req.Id)
	}

	// Parse update mask
	updateName := false
	updateDescription := false
	updateTags := false
	updateCategories := false

	if req.UpdateMask == nil || len(req.UpdateMask.Paths) == 0 {
		// If no update mask is provided, update all fields
		updateName = true
		updateDescription = true
		updateTags = true
		updateCategories = true
	} else {
		for _, path := range req.UpdateMask.Paths {
			switch path {
			case "name":
				updateName = true
			case "description":
				updateDescription = true
			case "tag_ids":
				updateTags = true
			case "category_ids":
				updateCategories = true
			}
		}
	}

	// Update exercise fields
	if updateName || updateDescription {
		sql := "UPDATE exercises SET"
		params := []interface{}{}
		needsComma := false

		if updateName {
			if req.Exercise.Name == "" {
				return nil, status.Error(codes.InvalidArgument, "exercise name cannot be empty")
			}
			sql += " name = ?"
			params = append(params, req.Exercise.Name)
			needsComma = true
		}

		if updateDescription {
			if needsComma {
				sql += ","
			}
			sql += " description = ?"
			params = append(params, req.Exercise.Description)
		}

		sql += " WHERE id = ?"
		params = append(params, req.Id)

		_, err = tx.ExecContext(ctx, sql, params...)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to update exercise: %v", err)
		}
	}

	// Update tags if requested
	if updateTags {
		// Delete existing tag associations
		_, err = tx.ExecContext(ctx, "DELETE FROM exercise_tags WHERE exercise_id = ?", req.Id)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to remove existing tag associations: %v", err)
		}

		// Add new tag associations
		for _, tagID := range req.Exercise.TagIds {
			// Check if tag exists
			var exists bool
			err := tx.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM tags WHERE id = ?)", tagID).Scan(&exists)
			if err != nil {
				return nil, status.Errorf(codes.Internal, "failed to check tag existence: %v", err)
			}
			if !exists {
				return nil, status.Errorf(codes.NotFound, "tag with ID %d not found", tagID)
			}

			// Create association
			_, err = tx.ExecContext(
				ctx,
				"INSERT INTO exercise_tags (exercise_id, tag_id) VALUES (?, ?)",
				req.Id, tagID,
			)
			if err != nil {
				return nil, status.Errorf(codes.Internal, "failed to associate exercise with tag: %v", err)
			}
		}
	}

	// Update categories if requested
	if updateCategories {
		// Delete existing category associations
		_, err = tx.ExecContext(ctx, "DELETE FROM exercise_categories WHERE exercise_id = ?", req.Id)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to remove existing category associations: %v", err)
		}

		// Add new category associations
		for _, categoryID := range req.Exercise.CategoryIds {
			// Check if category exists
			var exists bool
			err := tx.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM categories WHERE id = ?)", categoryID).Scan(&exists)
			if err != nil {
				return nil, status.Errorf(codes.Internal, "failed to check category existence: %v", err)
			}
			if !exists {
				return nil, status.Errorf(codes.NotFound, "category with ID %d not found", categoryID)
			}

			// Create association
			_, err = tx.ExecContext(
				ctx,
				"INSERT INTO exercise_categories (exercise_id, category_id) VALUES (?, ?)",
				req.Id, categoryID,
			)
			if err != nil {
				return nil, status.Errorf(codes.Internal, "failed to associate exercise with category: %v", err)
			}
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to commit transaction: %v", err)
	}

	// Retrieve the updated exercise
	return h.GetExercise(ctx, &pb.GetExerciseRequest{Id: req.Id})
}

// DeleteExercise deletes an exercise
func (h *ExerciseHandler) DeleteExercise(ctx context.Context, req *pb.DeleteExerciseRequest) (*emptypb.Empty, error) {
	if req.Id <= 0 {
		return nil, status.Error(codes.InvalidArgument, "invalid exercise ID")
	}

	// Check if the exercise exists
	var exists bool
	err := h.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM exercises WHERE id = ?)", req.Id).Scan(&exists)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to check exercise existence: %v", err)
	}
	if !exists {
		return nil, status.Errorf(codes.NotFound, "exercise with ID %d not found", req.Id)
	}

	// Delete the exercise (associated records in junction tables will be deleted by ON DELETE CASCADE)
	_, err = h.db.ExecContext(ctx, "DELETE FROM exercises WHERE id = ?", req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete exercise: %v", err)
	}

	return &emptypb.Empty{}, nil
}

// AddExerciseImage adds an image to an exercise
func (h *ExerciseHandler) AddExerciseImage(ctx context.Context, req *pb.AddExerciseImageRequest) (*pb.ExerciseImage, error) {
	if req.ExerciseId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "invalid exercise ID")
	}

	if len(req.ImageData) == 0 {
		return nil, status.Error(codes.InvalidArgument, "image data is required")
	}

	// Check if exercise exists
	var exists bool
	err := h.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM exercises WHERE id = ?)", req.ExerciseId).Scan(&exists)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to check exercise existence: %v", err)
	}
	if !exists {
		return nil, status.Errorf(codes.NotFound, "exercise with ID %d not found", req.ExerciseId)
	}

	// Insert the image
	result, err := h.db.ExecContext(
		ctx,
		`INSERT INTO exercise_images (exercise_id, image_data, filename, mime_type, description) 
         VALUES (?, ?, ?, ?, ?)`,
		req.ExerciseId, req.ImageData, req.Filename, req.MimeType, req.Description,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to add image to exercise: %v", err)
	}

	// Get the image ID
	imageID, err := result.LastInsertId()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get image ID: %v", err)
	}

	// Get image creation time
	var createdAt time.Time
	err = h.db.QueryRowContext(
		ctx,
		"SELECT created_at FROM exercise_images WHERE id = ?",
		imageID,
	).Scan(&createdAt)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get image creation time: %v", err)
	}

	// Return the created image
	return &pb.ExerciseImage{
		Id:          int32(imageID),
		ExerciseId:  req.ExerciseId,
		ImageData:   req.ImageData,
		Filename:    req.Filename,
		MimeType:    req.MimeType,
		Description: req.Description,
		CreatedAt:   timestamppb.New(createdAt),
	}, nil
}

// DeleteExerciseImage deletes an image from an exercise
func (h *ExerciseHandler) DeleteExerciseImage(ctx context.Context, req *pb.DeleteExerciseImageRequest) (*emptypb.Empty, error) {
	if req.Id <= 0 {
		return nil, status.Error(codes.InvalidArgument, "invalid image ID")
	}

	// Check if the image exists
	var exists bool
	err := h.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM exercise_images WHERE id = ?)", req.Id).Scan(&exists)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to check image existence: %v", err)
	}
	if !exists {
		return nil, status.Errorf(codes.NotFound, "image with ID %d not found", req.Id)
	}

	// Delete the image
	_, err = h.db.ExecContext(ctx, "DELETE FROM exercise_images WHERE id = ?", req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete image: %v", err)
	}

	return &emptypb.Empty{}, nil
}

// AddExerciseLink adds a link to an exercise
func (h *ExerciseHandler) AddExerciseLink(ctx context.Context, req *pb.AddExerciseLinkRequest) (*pb.ExerciseLink, error) {
	if req.ExerciseId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "invalid exercise ID")
	}

	if req.Url == "" {
		return nil, status.Error(codes.InvalidArgument, "url is required")
	}

	// Check if exercise exists
	var exists bool
	err := h.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM exercises WHERE id = ?)", req.ExerciseId).Scan(&exists)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to check exercise existence: %v", err)
	}
	if !exists {
		return nil, status.Errorf(codes.NotFound, "exercise with ID %d not found", req.ExerciseId)
	}

	// Insert the link
	result, err := h.db.ExecContext(
		ctx,
		`INSERT INTO exercise_links (exercise_id, url, description) VALUES (?, ?, ?)`,
		req.ExerciseId, req.Url, req.Description,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to add link to exercise: %v", err)
	}

	// Get the link ID
	linkID, err := result.LastInsertId()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get link ID: %v", err)
	}

	// Get link creation time
	var createdAt time.Time
	err = h.db.QueryRowContext(
		ctx,
		"SELECT created_at FROM exercise_links WHERE id = ?",
		linkID,
	).Scan(&createdAt)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get link creation time: %v", err)
	}

	// Return the created link
	return &pb.ExerciseLink{
		Id:          int32(linkID),
		ExerciseId:  req.ExerciseId,
		Url:         req.Url,
		Description: req.Description,
		CreatedAt:   timestamppb.New(createdAt),
	}, nil
}

// DeleteExerciseLink deletes a link from an exercise
func (h *ExerciseHandler) DeleteExerciseLink(ctx context.Context, req *pb.DeleteExerciseLinkRequest) (*emptypb.Empty, error) {
	if req.Id <= 0 {
		return nil, status.Error(codes.InvalidArgument, "invalid link ID")
	}

	// Check if the link exists
	var exists bool
	err := h.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM exercise_links WHERE id = ?)", req.Id).Scan(&exists)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to check link existence: %v", err)
	}
	if !exists {
		return nil, status.Errorf(codes.NotFound, "link with ID %d not found", req.Id)
	}

	// Delete the link
	_, err = h.db.ExecContext(ctx, "DELETE FROM exercise_links WHERE id = ?", req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete link: %v", err)
	}

	return &emptypb.Empty{}, nil
}

// GetExerciseStats retrieves statistics for an exercise
func (h *ExerciseHandler) GetExerciseStats(ctx context.Context, req *pb.GetExerciseStatsRequest) (*pb.ExerciseStats, error) {
	if req.ExerciseId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "invalid exercise ID")
	}

	// Check if exercise exists and get its name
	var exerciseName string
	err := h.db.QueryRowContext(
		ctx,
		"SELECT name FROM exercises WHERE id = ?",
		req.ExerciseId,
	).Scan(&exerciseName)
	if err == sql.ErrNoRows {
		return nil, status.Errorf(codes.NotFound, "exercise with ID %d not found", req.ExerciseId)
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to retrieve exercise: %v", err)
	}

	// Build date filter if provided
	dateFilter := ""
	dateParams := []interface{}{}
	if req.StartDate != nil {
		dateFilter += " AND start_time >= ?"
		dateParams = append(dateParams, req.StartDate.AsTime())
	}
	if req.EndDate != nil {
		dateFilter += " AND end_time <= ?"
		dateParams = append(dateParams, req.EndDate.AsTime())
	}

	// Get practice count
	var practiceCount int32
	practiceCountQuery := `
		SELECT COUNT(*) 
		FROM session_exercises 
		WHERE exercise_id = ?` + dateFilter
	params := append([]interface{}{req.ExerciseId}, dateParams...)
	err = h.db.QueryRowContext(ctx, practiceCountQuery, params...).Scan(&practiceCount)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get practice count: %v", err)
	}

	// Get total duration
	var totalDurationSeconds int32
	durationQuery := `
		SELECT COALESCE(SUM(strftime('%s', end_time) - strftime('%s', start_time)), 0) 
		FROM session_exercises 
		WHERE exercise_id = ?` + dateFilter
	err = h.db.QueryRowContext(ctx, durationQuery, params...).Scan(&totalDurationSeconds)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get total duration: %v", err)
	}

	// Get BPM statistics and ratings
	var maxBPM, minBPM int32
	var avgBPM, avgRating float64

	if practiceCount == 0 {
		maxBPM = 0
		minBPM = 0
		avgBPM = 0
		avgRating = 0
	} else {
		// Get max BPM
		maxBPMQuery := `
			SELECT COALESCE(MAX(bpm), 0) 
			FROM session_exercises 
			WHERE exercise_id = ? AND bpm IS NOT NULL` + dateFilter
		err = h.db.QueryRowContext(ctx, maxBPMQuery, params...).Scan(&maxBPM)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to get max BPM: %v", err)
		}

		// Get min BPM
		minBPMQuery := `
			SELECT COALESCE(MIN(bpm), 0) 
			FROM session_exercises 
			WHERE exercise_id = ? AND bpm IS NOT NULL` + dateFilter
		err = h.db.QueryRowContext(ctx, minBPMQuery, params...).Scan(&minBPM)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to get min BPM: %v", err)
		}

		// Get avg BPM
		avgBPMQuery := `
			SELECT COALESCE(AVG(bpm), 0) 
			FROM session_exercises 
			WHERE exercise_id = ? AND bpm IS NOT NULL` + dateFilter
		err = h.db.QueryRowContext(ctx, avgBPMQuery, params...).Scan(&avgBPM)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to get avg BPM: %v", err)
		}

		// Get avg rating from exercise_history
		historyParams := append([]interface{}{req.ExerciseId}, dateParams...)
		avgRatingQuery := `
			SELECT COALESCE(AVG(rating), 0) 
			FROM exercise_history 
			WHERE exercise_id = ? AND rating IS NOT NULL` + dateFilter
		err = h.db.QueryRowContext(ctx, avgRatingQuery, historyParams...).Scan(&avgRating)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to get avg rating: %v", err)
		}
	}

	// Get BPM progress over time
	bpmProgressQuery := `
		SELECT start_time, bpm 
		FROM session_exercises 
		WHERE exercise_id = ? AND bpm IS NOT NULL` + dateFilter + `
		ORDER BY start_time
	`
	bpmRows, err := h.db.QueryContext(ctx, bpmProgressQuery, params...)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get BPM progress: %v", err)
	}
	defer bpmRows.Close()

	var bpmProgress []*pb.BpmProgressPoint
	for bpmRows.Next() {
		var point pb.BpmProgressPoint
		var startTime time.Time
		if err := bpmRows.Scan(&startTime, &point.Bpm); err != nil {
			return nil, status.Errorf(codes.Internal, "failed to parse BPM progress: %v", err)
		}
		point.Date = timestamppb.New(startTime)
		bpmProgress = append(bpmProgress, &point)
	}
	if err = bpmRows.Err(); err != nil {
		return nil, status.Errorf(codes.Internal, "error reading BPM progress: %v", err)
	}

	// Return the statistics
	return &pb.ExerciseStats{
		ExerciseId:                   req.ExerciseId,
		ExerciseName:                 exerciseName,
		PracticeCount:                practiceCount,
		TotalPracticeDurationSeconds: totalDurationSeconds,
		AvgRating:                    avgRating,
		MaxBpm:                       maxBPM,
		MinBpm:                       minBPM,
		AvgBpm:                       avgBPM,
		BpmProgress:                  bpmProgress,
	}, nil
}
