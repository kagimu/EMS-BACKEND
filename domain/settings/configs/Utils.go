package configs

type Permission map[string][]string

func ReturnConfigs() []Permission {

	return []Permission{
		{"USERS": []string{"create", "read", "update", "delete"}},
		{"PERMISSIONS": []string{"create", "read", "update", "delete"}},
		{"INTERVIEWS": []string{"create", "read", "update", "delete"}},
		{"STUDENTS": []string{"create", "view", "update", "delete"}},
		{"SUBJECTS": []string{"create", "view", "update", "delete"}},
		{"LEARNING_AREAS": []string{"create", "view", "update", "delete"}},
		{"SUBJECT_CLUSTERS": []string{"create", "view", "update", "delete"}},
		{"INTERVIEW_GRADES": []string{"create", "view", "update", "delete"}},
		{"STUDENT_GRADES": []string{"create", "view", "update", "delete"}},
	}
	//response := make(map[string][]string)
	//response["USERS"] = []string{"create", "read", "update", "delete"}
	//response["PERMISSIONS"] = []string{"create", "read", "update", "delete"}
	//response["INTERVIEWS"] = []string{"create", "read", "update", "delete"}
	//return response

}
