package bot

import (
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/jaimeteb/chatto/clf"
	cmn "github.com/jaimeteb/chatto/common"
	"github.com/jaimeteb/chatto/ext"
	"github.com/jaimeteb/chatto/fsm"
	"github.com/spf13/viper"
)

// Bot models a bot with a Classifier and an FSM
type Bot struct {
	Name       string
	Machines   fsm.StoreFSM
	Domain     fsm.Domain
	Classifier clf.Classifier
	Extension  ext.Extension
	Clients    Clients
}

// Prediction models a classifier prediction and its orignal string
type Prediction struct {
	Original    string  `json:"original"`
	Predicted   string  `json:"predicted"`
	Probability float64 `json:"probability"`
}

// Config struct models the bot.yml configuration file
type Config struct {
	Name       string               `mapstructure:"bot_name"`
	Extensions ext.ExtensionsConfig `mapstructure:"extensions"`
	Store      fsm.StoreConfig      `mapstructure:"store"`
}

// Answer takes a user input and executes a transition on the FSM if possible
func (b Bot) Answer(mess cmn.Message) interface{} {
	if !b.Machines.Exists(mess.Sender) {
		b.Machines.Set(
			mess.Sender,
			&fsm.FSM{
				State: 0,
				Slots: make(map[string]string),
			},
		)
	}

	inputMessage := mess.Text
	cmd, _ := b.Classifier.Predict(inputMessage)

	m := b.Machines.Get(mess.Sender)
	resp, runExt := m.ExecuteCmd(cmd, inputMessage, b.Domain)
	if runExt != "" && b.Extension != nil {
		resp = b.Extension.RunExtFunc(mess.Sender, runExt, inputMessage, b.Domain, m)
	}
	b.Machines.Set(mess.Sender, m)

	return resp
}

// LoadBotConfig loads bot configuration from bot.yml
func LoadBotConfig(path *string) Config {
	config := viper.New()
	config.SetConfigName("bot")
	config.AddConfigPath(*path)
	config.AutomaticEnv()
	replacer := strings.NewReplacer(".", "_")
	config.SetEnvKeyReplacer(replacer)

	if err := config.ReadInConfig(); err != nil {
		switch err.(type) {
		case viper.ConfigFileNotFoundError:
			log.Warn("File bot.yml not found, using default values")
		default:
			log.Warn(err)
		}
		return Config{}
	}

	var bc Config
	config.Unmarshal(&bc)

	return bc
}

// LoadName loads the bot name from the configuration file
func LoadName(bcName string) (name string) {
	name = "botto"
	if bcName != "" {
		name = bcName
	}
	log.Infof("My name is '%v'\n", name)
	return
}

// LoadBot loads all configurations and returns a Bot
func LoadBot(path *string) Bot {
	bc := LoadBotConfig(path)

	// Load Name
	name := LoadName(bc.Name)
	// Load Domain
	domain := fsm.Create(path)
	// Load Classifier
	classifier := clf.Create(path)
	// Load Extensions
	extension := ext.LoadExtensions(bc.Extensions)
	// Load clients
	clients := LoadClients(path)
	// Load Store
	machines := fsm.LoadStore(bc.Store)

	return Bot{name, machines, domain, classifier, extension, clients}
}

// LOGO for Chatto
const LOGO = `
                           *******                          
                  *************************                 
             *********                *********             
          *******                           ******.         
        *****                                  ******       
      *****                                       *****     
    *****                                           *****   
   ****                                              .****  
  ****         ********,             *********         **** 
 ****       .******.******         ******.******       .****
 ****       ****       ****       ****       ****       ****
****                                                    ****
****                                                     ***
****                                                    ****
 ****                                                   ****
 ****                  ****       ****                 ****.
  ****                 *****     ****                  **** 
   ****.                 ***********                 *****  
    *****                                           ****    
      *****                                       *****     
        ******                                 ******       
          .******                          .******          
              *********               *********             
                  .***********************.                 
`
