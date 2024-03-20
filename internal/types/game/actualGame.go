package game







func HideCard(t PokerTable, eID int) PokerTable {
	res := t

	ps := []PokerPlayer{}

	for _, p := range t.Players {

		// if p.Action == "fold" {
		// 			newP := PokerPlayer{}

		// 	newP.ID = p.ID
		// 	newP.Username = p.Username
		// 	newP.Chips = p.Chips
		// 	newP.Action = p.Action
		// 	newP.Status = p.Status
		// 	newP.Img = p.Img
		// 	newP.TimeRemains = p.TimeRemains
		// 	newP.TimeBank = p.TimeBank
		// 	newP.Hand = p.Hand
		// 	newP.Bet = p.Bet
		// 	newP.Combination = p.Combination
		// 	newP.IsDealer = p.IsDealer
		// 	newP.NextAction = ""
		// 	ps = append(ps, newP)

		// 	continue
		// }
		
		if eID == p.ID {
			ps = append(ps, p)
		} else {
			newP := PokerPlayer{}

			newP.ID = p.ID
			newP.Username = p.Username
			newP.Chips = p.Chips
			newP.Action = p.Action
			newP.Status = p.Status
			newP.Img = p.Img
			newP.TimeRemains = p.TimeRemains
			newP.TimeBank = p.TimeBank
			newP.Hand = []string{"", ""}
			newP.Bet = p.Bet
			newP.Combination = p.Combination
			newP.IsDealer = p.IsDealer
			newP.NextAction = ""
			ps = append(ps, newP)
			
		}

	}

	res.Players = ps

	return res
}