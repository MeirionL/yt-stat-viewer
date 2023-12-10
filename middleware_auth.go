package main

// type authedHandler func(http.ResponseWriter, *http.Request, database.User)

// func (cfg *apiConfig) middlewareAuth(handler authedHandler) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		apiKey, err := auth.GetAPIKey(r.Header)
// 		if err != nil {
// 			RespondWithError(w, 403, fmt.Sprintf("auth error: %v", err))
// 			return
// 		}
// 		user, err := cfg.DB.GetUserByAPIKey(r.Context(), apiKey)
// 		if err != nil {
// 			RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("couldn't get users by api key: %v", err))
// 			return
// 		}

// 		handler(w, r, user)
// 	}
// }
