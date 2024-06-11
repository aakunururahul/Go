package main

func getNewUpdates(allUpdates []Message, createdAt int64) []Message {
	var updates []Message
	for _, update := range allUpdates {
		if update.CreatedAt > createdAt {
			updates = append(updates, update)
		}
	}
	return updates
}
