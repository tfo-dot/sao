package discord

import (
	"sao/player/inventory"
	"sao/world/party"

	"github.com/disgoorg/disgo/discord"
)

var PathToString = map[inventory.SkillPath]string{
	inventory.PathControl:   "Kontrola",
	inventory.PathDamage:    "Obrażenia",
	inventory.PathEndurance: "Wytrzymałość",
	inventory.PathMobility:  "Mobilność",
}

var RoleToString = map[party.PartyRole]string{
	party.Leader:  "Lider",
	party.DPS:     "DPS",
	party.Support: "Support",
	party.Tank:    "Tank",
}

var (
	DISCORD_COMMANDS = []discord.ApplicationCommandCreate{
		discord.SlashCommandCreate{
			Name:        "create",
			Description: "Stwórz postać",
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionString{
					Name:        "nazwa",
					Description: "Nazwa postaci",
					Required:    true,
				},
				discord.ApplicationCommandOptionUser{
					Name:        "gracz",
					Description: "Gracz",
					Required:    true,
				},
			},
		},
		discord.SlashCommandCreate{
			Name:        "ruch",
			Description: "Przenieś się do innej lokacji",
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionString{
					Autocomplete: true,
					Name:         "nazwa",
					Description:  "Nazwa lokacji",
					Required:     true,
				},
			},
		},
		discord.SlashCommandCreate{
			Name:        "tp",
			Description: "Teleportuj się na inne piętro",
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionString{
					Name:         "nazwa",
					Description:  "Nazwa piętra",
					Required:     true,
					Autocomplete: true,
				},
			},
		},
		discord.SlashCommandCreate{
			Name:        "info",
			Description: "Informacje o postaci",
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionUser{
					Name:        "gracz",
					Description: "Gracz",
					Required:    false,
				},
			},
		},
		discord.SlashCommandCreate{
			Name:        "skill",
			Description: "Zarządzaj umiejętnościami",
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionSubCommand{
					Name:        "pokaż",
					Description: "Pokaż umiejętności",
					Options: []discord.ApplicationCommandOption{
						discord.ApplicationCommandOptionUser{
							Name:        "gracz",
							Description: "Gracz",
							Required:    false,
						},
					},
				},
				discord.ApplicationCommandOptionSubCommand{
					Name:        "odblokuj",
					Description: "Odblokuj umiejętność",
					Options: []discord.ApplicationCommandOption{
						discord.ApplicationCommandOptionInt{
							Name:        "lvl",
							Description: "Umiejętność którego chcesz odblokować",
							Required:    true,
						},
					},
				},
				discord.ApplicationCommandOptionSubCommand{
					Name:        "ulepsz",
					Description: "Ulepsz umiejętność",
				},
			},
		},
		discord.SlashCommandCreate{
			Name:        "plecak",
			Description: "Zarządzaj ekwipunkiem",
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionSubCommand{
					Name:        "pokaż",
					Description: "Pokaż ekwipunek",
				},
			},
		},
		discord.SlashCommandCreate{
			Name:        "szukaj",
			Description: "Szukaj zajęcia",
		},
		discord.SlashCommandCreate{
			Name:        "party",
			Description: "Zarządzaj party",
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionSubCommand{
					Name:        "pokaż",
					Description: "Pokaż party",
				},
				discord.ApplicationCommandOptionSubCommand{
					Name:        "zapros",
					Description: "Zaproś do party",
					Options: []discord.ApplicationCommandOption{
						discord.ApplicationCommandOptionUser{
							Name:        "gracz",
							Description: "Gracz",
							Required:    true,
						},
					},
				},
				discord.ApplicationCommandOptionSubCommand{
					Name:        "wyrzuć",
					Description: "Wyrzuć z party",
					Options: []discord.ApplicationCommandOption{
						discord.ApplicationCommandOptionUser{
							Name:        "gracz",
							Description: "Gracz",
							Required:    true,
						},
					},
				},
				discord.ApplicationCommandOptionSubCommand{
					Name:        "opuść",
					Description: "Opuść party",
				},
				discord.ApplicationCommandOptionSubCommand{
					Name:        "zmień",
					Description: "Zmień rolę",
					Options: []discord.ApplicationCommandOption{
						discord.ApplicationCommandOptionUser{
							Name:        "gracz",
							Description: "Gracz",
							Required:    true,
						},
						discord.ApplicationCommandOptionString{
							Name:        "rola",
							Description: "Rola",
							Required:    true,
							Choices: []discord.ApplicationCommandOptionChoiceString{
								{
									Name:  "Lider",
									Value: "Lider",
								},
								{
									Name:  "DPS",
									Value: "DPS",
								},
								{
									Name:  "Support",
									Value: "Support",
								},
								{
									Name:  "Tank",
									Value: "Tank",
								},
							},
						},
					},
				},
				discord.ApplicationCommandOptionSubCommand{
					Name:        "rozwiąż",
					Description: "Rozwiąż party",
				},
			},
		},
		discord.SlashCommandCreate{
			Name:        "ratuj",
			Description: "Ratuj gracza",
		},
		discord.SlashCommandCreate{
			Name:        "speedrun",
			Description: "Czas przyśpieszyć!",
		},
	}
)

func isAdmin(member *discord.ResolvedMember) bool {
	return member.Permissions.Has(discord.PermissionAdministrator)
}

func MessageContent(content string, ephemeral bool) discord.MessageCreate {
	return discord.NewMessageCreateBuilder().SetContent(content).SetEphemeral(ephemeral).Build()
}
