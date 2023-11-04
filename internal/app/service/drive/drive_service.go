package drive

import (
	"context"
	"errors"
	"io"
	"log"
	"strings"

	domainDrive "github.com/onepiece010938/Line2GoogleDriveBot/internal/domain/drive"
	"golang.org/x/oauth2"
)

// For example
func (dr *GoogleDriveService) ListFiles(ctx context.Context, lineID string) (map[string]string, error) {
	// token改成去db取
	dToken, err := dr.driveServiceDynamodb.GetGoogleOAuthToken(lineID)
	// tok, err := dr.driveServiceGoogleOA.UserOAuthToken(authCode)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	tok := oauth2.Token{
		AccessToken:  dToken.AccessToken,
		TokenType:    dToken.TokenType,
		RefreshToken: dToken.RefreshToken,
		Expiry:       dToken.Expiry,
	}
	d, err := dr.driveServiceGoogleOA.NewGoogleDrive(ctx, &tok)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	result, err := d.ListFiles(10)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return result, nil
}

// For example
func (dr *GoogleDriveService) ListMyDriveFolders(ctx context.Context, lineID string) (map[string]string, error) {

	dToken, err := dr.driveServiceDynamodb.GetGoogleOAuthToken(lineID)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	tok := oauth2.Token{
		AccessToken:  dToken.AccessToken,
		TokenType:    dToken.TokenType,
		RefreshToken: dToken.RefreshToken,
		Expiry:       dToken.Expiry,
	}
	d, err := dr.driveServiceGoogleOA.NewGoogleDrive(ctx, &tok)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	result, err := d.ListMyDriveFolders()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return result, nil
}

// For example
func (dr *GoogleDriveService) ListSharedFolders(ctx context.Context, lineID string) (map[string]string, error) {

	dToken, err := dr.driveServiceDynamodb.GetGoogleOAuthToken(lineID)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	tok := oauth2.Token{
		AccessToken:  dToken.AccessToken,
		TokenType:    dToken.TokenType,
		RefreshToken: dToken.RefreshToken,
		Expiry:       dToken.Expiry,
	}
	d, err := dr.driveServiceGoogleOA.NewGoogleDrive(ctx, &tok)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	result, err := d.ListSharedFolders()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return result, nil
}

// For example
func (dr *GoogleDriveService) TestFolderCarousel(ctx context.Context, lineID string) (*domainDrive.FolderCarousel, error) {
	dToken, err := dr.driveServiceDynamodb.GetGoogleOAuthToken(lineID)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	tok := oauth2.Token{
		AccessToken:  dToken.AccessToken,
		TokenType:    dToken.TokenType,
		RefreshToken: dToken.RefreshToken,
		Expiry:       dToken.Expiry,
	}
	d, err := dr.driveServiceGoogleOA.NewGoogleDrive(ctx, &tok)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	myDrive, err := d.ListMyDriveFolders()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var params domainDrive.NewFolderCarouselParam

	for folderID, name := range myDrive {
		insideFolderM, err := d.ListFolderByID(folderID)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		fileM, err := d.ListFilesByID(folderID)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		param := domainDrive.NewFolderBubbleParam{
			Type:          "我的雲端硬碟",
			Name:          name,
			Path:          "/我的雲端硬碟/" + name,
			ID:            folderID,
			InsideFolderM: insideFolderM,
			FileM:         fileM,
		}
		params.BubbleParams = append(params.BubbleParams, param)
	}

	result, err := d.FindFolderPathByID("1kpLZfvk9XmSr4xtDvczAqYHIF8P3r8bk")
	log.Println("RES路徑:\n", result)
	result2, _ := d.FindFolderPathByID("1E9Gwyrwt4KMJ4NRk0FYg6dDH-6qwdaWs")
	log.Println("RES2路徑:\n", result2)
	// MinsideFolderM, err := d.ListFolderByID("1E9Gwyrwt4KMJ4NRk0FYg6dDH-6qwdaWs")
	// log.Println("MinsideFolderM:\n", MinsideFolderM)

	// 1EHHMlqSLrG7Q-1AetO1qa6yuC1PLSALP TestShared2
	result3, _ := d.FindFolderPathByID("1EHHMlqSLrG7Q-1AetO1qa6yuC1PLSALP")
	log.Println("RES3路徑:\n", result3)

	// if err != nil {
	// 	log.Println(err)
	// 	return nil, err
	// }
	// fileM, err := d.ListFiles(10)
	// if err != nil {
	// 	log.Println(err)
	// 	return nil, err
	// }

	// insideFolderM := map[string]string{
	// 	"001": "F1",
	// 	"002": "F2",
	// }
	// fileM := map[string]string{
	// 	"001": "file1",
	// 	"002": "file2",
	// }

	// var params domainDrive.NewFolderCarouselParam
	// params.BubbleParams = append(params.BubbleParams,
	// 	domainDrive.NewFolderBubbleParam{
	// 		Type:          "我的雲端硬碟",
	// 		Name:          "Folder1",
	// 		Path:          "/xx/xx",
	// 		ID:            "123",
	// 		InsideFolderM: insideFolderM,
	// 		FileM:         fileM,
	// 	},
	// 	domainDrive.NewFolderBubbleParam{
	// 		Type:          "我的雲端硬碟",
	// 		Name:          "Folder2",
	// 		Path:          "/yy/yy",
	// 		ID:            "1234",
	// 		InsideFolderM: insideFolderM,
	// 		FileM:         fileM,
	// 	},
	// )
	carousel := domainDrive.NewFolderCarousel(params)
	// return &carousel, nil
	return &carousel, err
}

func (dr *GoogleDriveService) ListFolderCarousel(ctx context.Context, lineID string, folderType domainDrive.FolderType) (*domainDrive.FolderCarousel, error) {
	dToken, err := dr.driveServiceDynamodb.GetGoogleOAuthToken(lineID)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	tok := oauth2.Token{
		AccessToken:  dToken.AccessToken,
		TokenType:    dToken.TokenType,
		RefreshToken: dToken.RefreshToken,
		Expiry:       dToken.Expiry,
	}
	d, err := dr.driveServiceGoogleOA.NewGoogleDrive(ctx, &tok)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var folderList map[string]string
	var folderTypeString string

	switch folderType {
	case domainDrive.PersonalFolder:
		folderList, err = d.ListMyDriveFolders()
		folderTypeString = "我的雲端硬碟"
	case domainDrive.SharedFolder:
		folderList, err = d.ListSharedFolders()
		folderTypeString = "與我共用"
	default:
		return nil, errors.New("unsupported folder type")
	}
	if err != nil {
		return nil, err
	}

	var params domainDrive.NewFolderCarouselParam
	// ROOTID, _ := d.GetRootID()
	// log.Println("RRRRRRRRR:", ROOTID)
	// rootID = ROOTID
	// // 做出當前目錄下的自己，才能看到資料夾下的檔案
	// currentFile, err := d.ListFilesByID(rootID)
	// if err != nil {
	// 	log.Println(err)
	// 	return nil, err
	// }
	// params.BubbleParams = append(params.BubbleParams, domainDrive.NewFolderBubbleParam{
	// 	Type:          folderTypeString,
	// 	Name:          folderTypeString,
	// 	Path:          "/" + folderTypeString,
	// 	ID:            rootID,
	// 	InsideFolderM: folderList,
	// 	FileM:         currentFile,
	// })

	for folderID, name := range folderList {
		path, err := d.FindFolderPathByID(folderID)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		insideFolderM, err := d.ListFolderByID(folderID)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		fileM, err := d.ListFilesByID(folderID)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		param := domainDrive.NewFolderBubbleParam{
			Type:          folderTypeString,
			Name:          name,
			Path:          path,
			ID:            folderID,
			InsideFolderM: insideFolderM,
			FileM:         fileM,
		}
		params.BubbleParams = append(params.BubbleParams, param)
	}

	carousel := domainDrive.NewFolderCarousel(params)

	return &carousel, err
}

func (dr *GoogleDriveService) ListSelectedFolderCarousel(ctx context.Context, lineID string, folderID string) (*domainDrive.FolderCarousel, error) {
	dToken, err := dr.driveServiceDynamodb.GetGoogleOAuthToken(lineID)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	tok := oauth2.Token{
		AccessToken:  dToken.AccessToken,
		TokenType:    dToken.TokenType,
		RefreshToken: dToken.RefreshToken,
		Expiry:       dToken.Expiry,
	}
	d, err := dr.driveServiceGoogleOA.NewGoogleDrive(ctx, &tok)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var params domainDrive.NewFolderCarouselParam

	// 做出當前目錄下的自己，才能看到資料夾下的檔案
	folderList, err := d.ListFolderByID(folderID)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	currentPath, err := d.FindFolderPathByID(folderID)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	currentFile, err := d.ListFilesByID(folderID)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	segments := strings.Split(currentPath, "/")
	currentName := segments[len(segments)-2]

	params.BubbleParams = append(params.BubbleParams, domainDrive.NewFolderBubbleParam{
		Type:          "打開資料夾",
		Name:          currentName,
		Path:          currentPath,
		ID:            folderID,
		InsideFolderM: folderList,
		FileM:         currentFile,
	})

	for folderID, name := range folderList {
		path, err := d.FindFolderPathByID(folderID)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		insideFolderM, err := d.ListFolderByID(folderID)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		fileM, err := d.ListFilesByID(folderID)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		param := domainDrive.NewFolderBubbleParam{
			Type:          "子資料夾",
			Name:          name,
			Path:          path,
			ID:            folderID,
			InsideFolderM: insideFolderM,
			FileM:         fileM,
		}
		params.BubbleParams = append(params.BubbleParams, param)
	}

	carousel := domainDrive.NewFolderCarousel(params)

	return &carousel, err
}

func (dr *GoogleDriveService) UploadFile(ctx context.Context, lineID string, fileName string, content io.ReadCloser) error {
	dToken, err := dr.driveServiceDynamodb.GetGoogleOAuthToken(lineID)

	if err != nil {
		log.Println(err)
		return err
	}

	tok := oauth2.Token{
		AccessToken:  dToken.AccessToken,
		TokenType:    dToken.TokenType,
		RefreshToken: dToken.RefreshToken,
		Expiry:       dToken.Expiry,
	}

	d, err := dr.driveServiceGoogleOA.NewGoogleDrive(ctx, &tok)
	if err != nil {
		log.Println(err)
		return err
	}

	file, err := domainDrive.SaveContent(content)
	if err != nil {
		log.Println(err)
		return err
	}

	log.Println("START Upload File To Drive")

	folderID := dToken.Info["upload_folder_id"].(string)

	err = d.UploadFile(folderID, fileName, file)
	if err != nil {
		log.Println("err:", err)
		return err
	}
	return nil

}

func (dr *GoogleDriveService) SetUploadPath(ctx context.Context, lineID string, folderID string) error {
	dToken, err := dr.driveServiceDynamodb.GetGoogleOAuthToken(lineID)
	if err != nil {
		log.Println(err)
		return err
	}
	dToken.PK = lineID
	dToken.Info = map[string]interface{}{
		"upload_folder_id": folderID,
	}

	_, err = dr.driveServiceDynamodb.TxUpdateGoogleOAuthToken(dToken)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (dr *GoogleDriveService) GetUploadPath(ctx context.Context, lineID string) (string, error) {
	dToken, err := dr.driveServiceDynamodb.GetGoogleOAuthToken(lineID)
	if err != nil {
		log.Println(err)
		return "", err
	}
	folderID := dToken.Info["upload_folder_id"]

	tok := oauth2.Token{
		AccessToken:  dToken.AccessToken,
		TokenType:    dToken.TokenType,
		RefreshToken: dToken.RefreshToken,
		Expiry:       dToken.Expiry,
	}

	d, err := dr.driveServiceGoogleOA.NewGoogleDrive(ctx, &tok)
	if err != nil {
		log.Println(err)
		return "", err
	}
	path, err := d.FindFolderPathByID(folderID.(string))
	if err != nil {
		log.Println(err)
		return "", err
	}
	return path, nil
}
