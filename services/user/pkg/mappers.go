package user

import "github.com/olo/litter3/models"

func (s *userService) mapSocietyToSocietyAnswSimple(societies []models.Society) []models.SocietyAnswSimple {
	var societiesAnsw []models.SocietyAnswSimple
	for _, society := range societies {
		societiesAnsw = append(societiesAnsw, models.SocietyAnswSimple{
			Id: society.Id,
			Name: society.Name,
			Avatar: society.Avatar,
			Description: society.Description,
			UsersNumb: len(society.Users),
			CreatedAt: society.CreatedAt,
		})
	}
	return societiesAnsw
}