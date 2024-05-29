package discord

import (
	"sao/types"
	"sao/world/party"

	"github.com/disgoorg/disgo/discord"
)

var PathToString = map[types.SkillPath]string{
	types.PathControl:   "Kontrola",
	types.PathDamage:    "Obrażenia",
	types.PathEndurance: "Wytrzymałość",
	types.PathMobility:  "Mobilność",
}

var RoleToString = map[party.PartyRole]string{
	party.DPS:     "DPS",
	party.Support: "Support",
	party.Tank:    "Tank",
	party.None:    "Brak",
}

var DISCORD_COMMANDS = []discord.ApplicationCommandCreate{
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
		Name:        "stwórz",
		Description: "Stwórz przedmiot",
		Options: []discord.ApplicationCommandOption{
			discord.ApplicationCommandOptionString{
				Name:         "nazwa",
				Description:  "Nazwa przedmiotu",
				Required:     true,
				Autocomplete: true,
			},
			discord.ApplicationCommandOptionInt{
				Name:        "ilość",
				Description: "Ilość przedmiotów",
			},
		},
	},
	discord.SlashCommandCreate{
		Name:        "furia",
		Description: "Zarządzaj furiami",
		Options: []discord.ApplicationCommandOption{
			discord.ApplicationCommandOptionSubCommand{
				Name:        "pokaż",
				Description: "Pokaż swoje furie",
			},
			discord.ApplicationCommandOptionSubCommand{
				Name:        "ulepsz",
				Description: "Ulepsz furie kolejny tier",
			},
		},
	},
	discord.SlashCommandCreate{
		Name:        "sklep",
		Description: "Zarządzaj sklepami",
		Options: []discord.ApplicationCommandOption{
			discord.ApplicationCommandOptionSubCommand{
				Name:        "pokaż",
				Description: "Pokazuje sklepy w danej lokalizacji",
			},
		},
	},
	discord.SlashCommandCreate{
		Name:        "turniej",
		Description: "Zarządzaj turniejami",
		Options: []discord.ApplicationCommandOption{
			discord.ApplicationCommandOptionSubCommand{
				Name:        "stwórz",
				Description: "Stwórz turniej",
				Options: []discord.ApplicationCommandOption{
					discord.ApplicationCommandOptionString{
						Name:        "nazwa",
						Description: "Nazwa turnieju",
						Required:    true,
					},
					discord.ApplicationCommandOptionInt{
						Name:        "typ",
						Description: "Typ turnieju",
						Required:    true,
						Choices: []discord.ApplicationCommandOptionChoiceInt{
							{
								Name:  "Zwykły",
								Value: 0,
							},
						},
					},
					discord.ApplicationCommandOptionInt{
						Name:        "max",
						Description: "Maksymalna ilość graczy",
					},
				},
			},
			discord.ApplicationCommandOptionSubCommand{
				Name:        "rozpocznij",
				Description: "Rozpocznij turniej",
				Options: []discord.ApplicationCommandOption{
					discord.ApplicationCommandOptionString{
						Name:         "nazwa",
						Description:  "Nazwa turnieju",
						Required:     true,
						Autocomplete: true,
					},
				},
			},
		},
	},
	discord.SlashCommandCreate{
		Name:        "handel",
		Description: "Zarządzaj handlem",
		Options: []discord.ApplicationCommandOption{
			discord.ApplicationCommandOptionSubCommand{
				Name:        "nowy",
				Description: "Zaproś do handlu",
				Options: []discord.ApplicationCommandOption{
					discord.ApplicationCommandOptionUser{
						Name:        "gracz",
						Description: "Kogo?",
						Required:    true,
					},
				},
			},
		},
	},
}

func isAdmin(member *discord.ResolvedMember) bool {
	return member.Permissions.Has(discord.PermissionAdministrator)
}

func MessageContent(content string, ephemeral bool) discord.MessageCreate {
	return discord.NewMessageCreateBuilder().SetContent(content).SetEphemeral(ephemeral).Build()
}

func MessageEmbed(embeds ...discord.Embed) discord.MessageCreate {
	return discord.NewMessageCreateBuilder().SetEmbeds(embeds...).Build()
}
