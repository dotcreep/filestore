package utils

type Success struct {
	Success bool   `json:"success" example:"true"`
	Result  string `json:"result" example:"message"`
	Message string `json:"message" example:"message"`
	Status  int    `json:"status" example:"200"`
	Error   string `json:"error" example:"null"`
}
type BadRequest struct {
	Success bool   `json:"success" example:"false"`
	Result  string `json:"result" example:"null"`
	Message string `json:"message" example:"message"`
	Status  int    `json:"status" example:"400"`
	Error   string `json:"error" example:"string"`
}
type InternalServerError struct {
	Success bool   `json:"success" example:"false"`
	Result  string `json:"result" example:"null"`
	Message string `json:"message" example:"message"`
	Status  int    `json:"status" example:"500"`
	Error   string `json:"error" example:"message"`
}
type Showfile struct {
	Success bool `json:"success" example:"true"`
	Result  struct {
		AAB struct {
			Hash struct {
				Index       int    `json:"index" example:"1"`
				Filename    string `json:"filename" example:"file.aab"`
				URL         string `json:"url" example:"/username/hash"`
				UploadAt    string `json:"upload_at" example:"2025-01-27T19:44:25.467738468+07:00"`
				LabelName   string `json:"label_name" example:"Example Apps"`
				Version     string `json:"version" example:"v1.0.0"`
				PackageName string `json:"package_name" example:"id.co.example.username"`
			} `json:"hash.aab"`
		} `json:"aab"`
		APK struct {
			Hash struct {
				Index       int    `json:"index" example:"1"`
				Filename    string `json:"filename" example:"file.apk"`
				URL         string `json:"url" example:"/username/hash"`
				UploadAt    string `json:"upload_at" example:"2025-01-27T19:44:25.467738468+07:00"`
				LabelName   string `json:"label_name" example:"Example Apps"`
				Version     string `json:"version" example:"v1.0.0"`
				PackageName string `json:"package_name" example:"id.co.example.username"`
			} `json:"hash.apk"`
		} `json:"apk"`
	} `json:"result"`
	Message string `json:"message" example:"message"`
	Status  int    `json:"status" example:"200"`
	Error   string `json:"error" example:"null"`
}
type FileData struct {
	Index       int    `json:"index" example:"1"`
	Filename    string `json:"filename" example:"file.apk"`
	URL         string `json:"url" example:"/username/hash"`
	UploadAt    string `json:"upload_at" example:"2025-01-27T19:44:25.467738468+07:00"`
	LabelName   string `json:"label_name" example:"Example Apps"`
	Version     string `json:"version" example:"v1.0.0"`
	PackageName string `json:"package_name" example:"id.co.example.username"`
}
