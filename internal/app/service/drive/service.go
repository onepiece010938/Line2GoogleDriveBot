package drive

import "context"

type GoogleDriveService struct {
	driveServiceGoogleOA DriveServiceGoogleOAuthI
	driveServiceDynamodb DriveServiceDynamodbI
}

type GoogleDriveServiceParam struct {
	DriveServiceGoogleOA DriveServiceGoogleOAuthI
	DriveServiceDynamodb DriveServiceDynamodbI
}

func NewGoogleDriveService(_ context.Context, param GoogleDriveServiceParam) *GoogleDriveService {
	return &GoogleDriveService{
		driveServiceGoogleOA: param.DriveServiceGoogleOA,
		driveServiceDynamodb: param.DriveServiceDynamodb,
	}
}
