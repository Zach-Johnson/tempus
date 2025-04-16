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

// CategoryHandler implements the CategoryService gRPC service
type CategoryHandler struct {
	pb.UnimplementedCategoryServiceServer
	db *sql.DB
}

// NewCategoryHandler creates a new CategoryHandler
func NewCategoryHandler(db *sql.DB) *CategoryHandler {
	return &CategoryHandler{db: db}
}

// CreateCategory creates a new category
func (h *CategoryHandler) CreateCategory(ctx context.Context, req *pb.CreateCategoryRequest) (*pb.Category, error) {
	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "category name is required")
	}

	var id int64
	var createdAt, updatedAt time.Time

	// Insert the category into the database
	query := "INSERT INTO categories (name, description) VALUES (?, ?)"
	result, err := h.db.ExecContext(ctx, query, req.Name, req.Description)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create category: %v", err)
	}

	// Get the generated ID
	id, err = result.LastInsertId()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get generated ID: %v", err)
	}

	// Fetch the created category to get timestamps
	err = h.db.QueryRowContext(
		ctx,
		"SELECT created_at, updated_at FROM categories WHERE id = ?",
		id,
	).Scan(&createdAt, &updatedAt)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to fetch created category: %v", err)
	}

	// Return the created category
	return &pb.Category{
		Id:          int32(id),
		Name:        req.Name,
		Description: req.Description,
		CreatedAt:   timestamppb.New(createdAt),
		UpdatedAt:   timestamppb.New(updatedAt),
	}, nil
}

// GetCategory retrieves a category by ID
func (h *CategoryHandler) GetCategory(ctx context.Context, req *pb.GetCategoryRequest) (*pb.Category, error) {
	if req.Id <= 0 {
		return nil, status.Error(codes.InvalidArgument, "invalid category ID")
	}

	var category pb.Category
	var createdAt, updatedAt time.Time

	// Query the database for the category
	err := h.db.QueryRowContext(
		ctx,
		"SELECT id, name, description, created_at, updated_at FROM categories WHERE id = ?",
		req.Id,
	).Scan(&category.Id, &category.Name, &category.Description, &createdAt, &updatedAt)

	if err == sql.ErrNoRows {
		return nil, status.Errorf(codes.NotFound, "category with ID %d not found", req.Id)
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to retrieve category: %v", err)
	}

	category.CreatedAt = timestamppb.New(createdAt)
	category.UpdatedAt = timestamppb.New(updatedAt)

	return &category, nil
}

// ListCategories lists all categories with pagination
func (h *CategoryHandler) ListCategories(ctx context.Context, req *pb.ListCategoriesRequest) (*pb.ListCategoriesResponse, error) {
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

	// Query total count
	var totalCount int32
	err := h.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM categories").Scan(&totalCount)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to count categories: %v", err)
	}

	// Query categories with pagination
	rows, err := h.db.QueryContext(
		ctx,
		"SELECT id, name, description, created_at, updated_at FROM categories ORDER BY name LIMIT ? OFFSET ?",
		pageSize+1, // Query one more to check if there are more pages
		offset,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list categories: %v", err)
	}
	defer rows.Close()

	// Parse the results
	categories := make([]*pb.Category, 0, pageSize)
	count := 0
	hasMorePages := false

	for rows.Next() {
		if count >= pageSize {
			hasMorePages = true
			break
		}

		var category pb.Category
		var createdAt, updatedAt time.Time

		err := rows.Scan(&category.Id, &category.Name, &category.Description, &createdAt, &updatedAt)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to parse category: %v", err)
		}

		category.CreatedAt = timestamppb.New(createdAt)
		category.UpdatedAt = timestamppb.New(updatedAt)
		categories = append(categories, &category)
		count++
	}

	if err = rows.Err(); err != nil {
		return nil, status.Errorf(codes.Internal, "error reading categories: %v", err)
	}

	// Calculate next page token
	nextPageToken := ""
	if hasMorePages {
		nextPageToken = strconv.Itoa(offset + pageSize)
	}

	return &pb.ListCategoriesResponse{
		Categories:    categories,
		NextPageToken: nextPageToken,
		TotalCount:    totalCount,
	}, nil
}

// UpdateCategory updates an existing category
func (h *CategoryHandler) UpdateCategory(ctx context.Context, req *pb.UpdateCategoryRequest) (*pb.Category, error) {
	if req.Id <= 0 {
		return nil, status.Error(codes.InvalidArgument, "invalid category ID")
	}

	if req.Category == nil {
		return nil, status.Error(codes.InvalidArgument, "category data is required")
	}

	// Check if the category exists
	var exists bool
	err := h.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM categories WHERE id = ?)", req.Id).Scan(&exists)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to check category existence: %v", err)
	}
	if !exists {
		return nil, status.Errorf(codes.NotFound, "category with ID %d not found", req.Id)
	}

	// Parse update mask
	updateName := false
	updateDescription := false

	if req.UpdateMask == nil || len(req.UpdateMask.Paths) == 0 {
		// If no update mask is provided, update all fields
		updateName = true
		updateDescription = true
	} else {
		for _, path := range req.UpdateMask.Paths {
			switch path {
			case "name":
				updateName = true
			case "description":
				updateDescription = true
			}
		}
	}

	// Build update SQL
	sql := "UPDATE categories SET"
	params := []any{}
	first := true

	if updateName {
		if req.Category.Name == "" {
			return nil, status.Error(codes.InvalidArgument, "category name cannot be empty")
		}
		sql += " name = ?"
		params = append(params, req.Category.Name)
		first = false
	}

	if updateDescription {
		if !first {
			sql += ","
		}
		sql += " description = ?"
		params = append(params, req.Category.Description)
	}

	if len(params) == 0 {
		return nil, status.Error(codes.InvalidArgument, "no fields to update")
	}

	sql += " WHERE id = ?"
	params = append(params, req.Id)

	// Execute the update
	_, err = h.db.ExecContext(ctx, sql, params...)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update category: %v", err)
	}

	// Retrieve the updated category
	return h.GetCategory(ctx, &pb.GetCategoryRequest{Id: req.Id})
}

// DeleteCategory deletes a category
func (h *CategoryHandler) DeleteCategory(ctx context.Context, req *pb.DeleteCategoryRequest) (*emptypb.Empty, error) {
	if req.Id <= 0 {
		return nil, status.Error(codes.InvalidArgument, "invalid category ID")
	}

	// Check if the category exists
	var exists bool
	err := h.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM categories WHERE id = ?)", req.Id).Scan(&exists)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to check category existence: %v", err)
	}
	if !exists {
		return nil, status.Errorf(codes.NotFound, "category with ID %d not found", req.Id)
	}

	// Delete the category
	_, err = h.db.ExecContext(ctx, "DELETE FROM categories WHERE id = ?", req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete category: %v", err)
	}

	return &emptypb.Empty{}, nil
}
