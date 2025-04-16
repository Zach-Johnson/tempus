package handlers

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"time"

	pb "github.com/Zach-Johnson/tempus/proto/api/v1/tempus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// TagService implements the TagService gRPC service
type TagService struct {
	pb.UnimplementedTagServiceServer
	db *sql.DB
}

// NewTagService creates a new TagService
func NewTagService(db *sql.DB) *TagService {
	return &TagService{db: db}
}

// CreateTag creates a new tag
func (s *TagService) CreateTag(ctx context.Context, req *pb.CreateTagRequest) (*pb.Tag, error) {
	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "tag name is required")
	}

	// Start a transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to start transaction: %v", err)
	}
	defer tx.Rollback() // Rollback if not committed

	// Insert the tag
	var id int32
	var createdAt time.Time

	err = tx.QueryRowContext(
		ctx,
		"INSERT INTO tags (name) VALUES (?) RETURNING id, created_at",
		req.Name,
	).Scan(&id, &createdAt)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create tag: %v", err)
	}

	// Associate tag with categories if provided
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

		// Create tag-category relationship
		_, err = tx.ExecContext(
			ctx,
			"INSERT INTO tag_categories (tag_id, category_id) VALUES (?, ?)",
			id, categoryID,
		)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to associate tag with category: %v", err)
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to commit transaction: %v", err)
	}

	// Return the created tag
	return &pb.Tag{
		Id:          id,
		Name:        req.Name,
		CreatedAt:   timestamppb.New(createdAt),
		CategoryIds: req.CategoryIds,
	}, nil
}

// GetTag retrieves a tag by ID
func (s *TagService) GetTag(ctx context.Context, req *pb.GetTagRequest) (*pb.Tag, error) {
	if req.Id <= 0 {
		return nil, status.Error(codes.InvalidArgument, "invalid tag ID")
	}

	// Check if tag exists
	var tag pb.Tag
	var createdAt time.Time

	err := s.db.QueryRowContext(
		ctx,
		"SELECT id, name, created_at FROM tags WHERE id = ?",
		req.Id,
	).Scan(&tag.Id, &tag.Name, &createdAt)

	if err == sql.ErrNoRows {
		return nil, status.Errorf(codes.NotFound, "tag with ID %d not found", req.Id)
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to retrieve tag: %v", err)
	}

	tag.CreatedAt = timestamppb.New(createdAt)

	// Get associated category IDs
	rows, err := s.db.QueryContext(
		ctx,
		"SELECT category_id FROM tag_categories WHERE tag_id = ?",
		req.Id,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to retrieve tag categories: %v", err)
	}
	defer rows.Close()

	var categoryIDs []int32
	for rows.Next() {
		var categoryID int32
		if err := rows.Scan(&categoryID); err != nil {
			return nil, status.Errorf(codes.Internal, "failed to parse category ID: %v", err)
		}
		categoryIDs = append(categoryIDs, categoryID)
	}

	if err = rows.Err(); err != nil {
		return nil, status.Errorf(codes.Internal, "error reading tag categories: %v", err)
	}

	tag.CategoryIds = categoryIDs

	return &tag, nil
}

// ListTags lists all tags with pagination and optional filtering
func (s *TagService) ListTags(ctx context.Context, req *pb.ListTagsRequest) (*pb.ListTagsResponse, error) {
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
	var countQuery, query string
	var countParams, queryParams []any

	if req.CategoryId > 0 {
		// Filter by category
		countQuery = `
            SELECT COUNT(DISTINCT t.id) 
            FROM tags t
            JOIN tag_categories tc ON t.id = tc.tag_id
            WHERE tc.category_id = ?
        `
		countParams = []any{req.CategoryId}

		query = `
            SELECT DISTINCT t.id, t.name, t.created_at
            FROM tags t
            JOIN tag_categories tc ON t.id = tc.tag_id
            WHERE tc.category_id = ?
            ORDER BY t.name
            LIMIT ? OFFSET ?
        `
		queryParams = []any{req.CategoryId, pageSize + 1, offset}
	} else {
		// No filter
		countQuery = "SELECT COUNT(*) FROM tags"
		countParams = []any{}

		query = `
            SELECT id, name, created_at 
            FROM tags 
            ORDER BY name
            LIMIT ? OFFSET ?
        `
		queryParams = []any{pageSize + 1, offset}
	}

	// Query total count
	var totalCount int32
	err := s.db.QueryRowContext(ctx, countQuery, countParams...).Scan(&totalCount)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to count tags: %v", err)
	}

	// Query tags with pagination
	rows, err := s.db.QueryContext(ctx, query, queryParams...)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list tags: %v", err)
	}
	defer rows.Close()

	// Parse the results
	tags := make([]*pb.Tag, 0, pageSize)
	count := 0
	hasMorePages := false
	var tagIDs []int32

	for rows.Next() {
		if count >= pageSize {
			hasMorePages = true
			break
		}

		var tag pb.Tag
		var createdAt time.Time

		err := rows.Scan(&tag.Id, &tag.Name, &createdAt)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to parse tag: %v", err)
		}

		tag.CreatedAt = timestamppb.New(createdAt)
		tags = append(tags, &tag)
		tagIDs = append(tagIDs, tag.Id)
		count++
	}

	if err = rows.Err(); err != nil {
		return nil, status.Errorf(codes.Internal, "error reading tags: %v", err)
	}

	// Get category IDs for all returned tags
	if len(tagIDs) > 0 {
		// Build query with placeholders for all tag IDs
		placeholders := ""
		params := []any{}
		for i, id := range tagIDs {
			if i > 0 {
				placeholders += ", "
			}
			placeholders += "?"
			params = append(params, id)
		}

		categoryRows, err := s.db.QueryContext(
			ctx,
			fmt.Sprintf("SELECT tag_id, category_id FROM tag_categories WHERE tag_id IN (%s)", placeholders),
			params...,
		)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to retrieve tag categories: %v", err)
		}
		defer categoryRows.Close()

		// Map tag_id to slice of category_ids
		tagToCategoryMap := make(map[int32][]int32)
		for categoryRows.Next() {
			var tagID, categoryID int32
			if err := categoryRows.Scan(&tagID, &categoryID); err != nil {
				return nil, status.Errorf(codes.Internal, "failed to parse tag category: %v", err)
			}
			tagToCategoryMap[tagID] = append(tagToCategoryMap[tagID], categoryID)
		}

		if err = categoryRows.Err(); err != nil {
			return nil, status.Errorf(codes.Internal, "error reading tag categories: %v", err)
		}

		// Add category IDs to tags
		for _, tag := range tags {
			tag.CategoryIds = tagToCategoryMap[tag.Id]
		}
	}

	// Calculate next page token
	nextPageToken := ""
	if hasMorePages {
		nextPageToken = strconv.Itoa(offset + pageSize)
	}

	return &pb.ListTagsResponse{
		Tags:          tags,
		NextPageToken: nextPageToken,
		TotalCount:    totalCount,
	}, nil
}

// UpdateTag updates an existing tag
func (s *TagService) UpdateTag(ctx context.Context, req *pb.UpdateTagRequest) (*pb.Tag, error) {
	if req.Id <= 0 {
		return nil, status.Error(codes.InvalidArgument, "invalid tag ID")
	}

	if req.Tag == nil {
		return nil, status.Error(codes.InvalidArgument, "tag data is required")
	}

	// Start a transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to start transaction: %v", err)
	}
	defer tx.Rollback() // Rollback if not committed

	// Check if the tag exists
	var exists bool
	err = tx.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM tags WHERE id = ?)", req.Id).Scan(&exists)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to check tag existence: %v", err)
	}
	if !exists {
		return nil, status.Errorf(codes.NotFound, "tag with ID %d not found", req.Id)
	}

	// Parse update mask
	updateName := false
	updateCategories := false

	if req.UpdateMask == nil || len(req.UpdateMask.Paths) == 0 {
		// If no update mask is provided, update all fields
		updateName = true
		updateCategories = true
	} else {
		for _, path := range req.UpdateMask.Paths {
			switch path {
			case "name":
				updateName = true
			case "category_ids":
				updateCategories = true
			}
		}
	}

	// Update name if requested
	if updateName {
		if req.Tag.Name == "" {
			return nil, status.Error(codes.InvalidArgument, "tag name cannot be empty")
		}

		_, err = tx.ExecContext(
			ctx,
			"UPDATE tags SET name = ? WHERE id = ?",
			req.Tag.Name, req.Id,
		)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to update tag name: %v", err)
		}
	}

	// Update category associations if requested
	if updateCategories {
		// Delete existing associations
		_, err = tx.ExecContext(
			ctx,
			"DELETE FROM tag_categories WHERE tag_id = ?",
			req.Id,
		)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to delete tag categories: %v", err)
		}

		// Add new associations
		for _, categoryID := range req.Tag.CategoryIds {
			// Check if category exists
			var exists bool
			err := tx.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM categories WHERE id = ?)", categoryID).Scan(&exists)
			if err != nil {
				return nil, status.Errorf(codes.Internal, "failed to check category existence: %v", err)
			}
			if !exists {
				return nil, status.Errorf(codes.NotFound, "category with ID %d not found", categoryID)
			}

			// Create tag-category relationship
			_, err = tx.ExecContext(
				ctx,
				"INSERT INTO tag_categories (tag_id, category_id) VALUES (?, ?)",
				req.Id, categoryID,
			)
			if err != nil {
				return nil, status.Errorf(codes.Internal, "failed to associate tag with category: %v", err)
			}
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to commit transaction: %v", err)
	}

	// Retrieve the updated tag
	return s.GetTag(ctx, &pb.GetTagRequest{Id: req.Id})
}

// DeleteTag deletes a tag
func (s *TagService) DeleteTag(ctx context.Context, req *pb.DeleteTagRequest) (*emptypb.Empty, error) {
	if req.Id <= 0 {
		return nil, status.Error(codes.InvalidArgument, "invalid tag ID")
	}

	// Check if the tag exists
	var exists bool
	err := s.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM tags WHERE id = ?)", req.Id).Scan(&exists)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to check tag existence: %v", err)
	}
	if !exists {
		return nil, status.Errorf(codes.NotFound, "tag with ID %d not found", req.Id)
	}

	// Delete the tag (associated records in junction tables will be deleted by ON DELETE CASCADE)
	_, err = s.db.ExecContext(ctx, "DELETE FROM tags WHERE id = ?", req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete tag: %v", err)
	}

	return &emptypb.Empty{}, nil
}
