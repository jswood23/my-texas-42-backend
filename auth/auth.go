package auth

//func getStore(c *gin.Context) (*pgstore.PGStore, error) {
//	// Fetch new store
//	database, err := services.GetDBSession()
//	if err != nil {
//		responseDatabaseError := make(map[string]string)
//		responseDatabaseError["message"] = "There was a problem connecting to the database."
//
//		c.JSON(http.StatusInternalServerError, responseDatabaseError)
//		c.Abort()
//		return nil, err
//	}
//
//	store, err := pgstore.NewPGStoreFromPool(database, []byte(services.GetAuthSecretKey()))
//	if err != nil {
//		responseAuthError := make(map[string]string)
//		responseAuthError["message"] = "There was a problem authenticating the user."
//
//		c.JSON(http.StatusInternalServerError, responseAuthError)
//		c.Abort()
//		return nil, err
//	}
//	store.Options.HttpOnly = true
//
//	return store, nil
//}
//
//func Authenticate(c *gin.Context) {
//	store, err := getStore(c)
//	defer store.Close()
//	if err != nil {
//		println(err.Error())
//		return
//	}
//
//	// Run a background goroutine to clean up expired sessions from the database
//	defer store.StopCleanup(store.Cleanup(time.Minute * 5))
//
//	// Get a session
//	session, err := store.Get(c.Request, "session")
//	if err == nil && len(session.ID) > 0 {
//		c.Next()
//	} else {
//		if err != nil {
//			println(err.Error())
//		}
//
//		responseUnauthorized := make(map[string]string)
//		responseUnauthorized["message"] = "Invalid session key."
//
//		c.JSON(http.StatusUnauthorized, responseUnauthorized)
//		c.Abort()
//		return
//	}
//}
//
//func NewSession(c *gin.Context) error {
//	store, err := getStore(c)
//	defer store.Close()
//	if err != nil {
//		return err
//	}
//
//	// Run a background goroutine to clean up expired sessions from the database
//	defer store.StopCleanup(store.Cleanup(time.Minute * 5))
//
//	// Start a new session
//	session, err := store.New(c.Request, "session")
//	if err != nil {
//		return err
//	}
//
//	err = session.Save(c.Request, c.Writer)
//	if err != nil {
//		return err
//	}
//
//	println("Session ID: " + session.ID)
//	return nil
//}
