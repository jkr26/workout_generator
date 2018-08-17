package main

import (
"io/ioutil"
"fmt"
"log"
"encoding/json"
"math/rand"
"time"
"strings"
"github.com/zmb3/spotify"
)

type workoutClass struct{
	Name string `json:"name"`
	Attrs map[string][]string `json:"attrs"`
}

type Move struct{
	Name string
	Attrs map[string]string
}


type allWorkouts struct{
	Workouts []workoutClass `json:"workouts"`
}

func readFile(path string) ([]byte, error) {
	b, err := ioutil.ReadFile(path)
    return b, err
}

func remove(s []string, i int) []string {
    s[len(s)-1], s[i] = s[i], s[len(s)-1]
    return s[:len(s)-1]
}

func select_option_and_return_workout(w []string, o []string) ([]string, []string) {
	rand.Seed(time.Now().Unix())
	move := rand.Intn(len(o))
	w = append(w, o[move])
	o = remove(o, move)
	return w, o
}

func gen_workout_skeleton() ([]string) {
	//Function in which we define the allowable structure of workouts.
	var second_move_opts []string
	var third_move_opts []string
	var fourth_move_opts []string
	var fifth_move_opts []string
	workout := []string {"warmup"}
	first_move_opts := []string {"rowing", "intervals", "circuit"}
	workout, second_move_opts = select_option_and_return_workout(workout, first_move_opts)
	second_move_opts = append(second_move_opts, "lbstrength")
	workout, third_move_opts = select_option_and_return_workout(workout, second_move_opts)
	workout, fourth_move_opts = select_option_and_return_workout(workout, third_move_opts)
	fourth_move_opts = append(fourth_move_opts, "circuit")
	fourth_move_opts = append(fourth_move_opts, "ubstrength")
	workout, fifth_move_opts = select_option_and_return_workout(workout, fourth_move_opts)
	workout, _ = select_option_and_return_workout(workout, fifth_move_opts)
	workout = append(workout, "finisher")
	return workout
}	

func grab_random_workout(w workoutClass) (Move) {
	rand.Seed(time.Now().Unix())
	n := w.Name
	strengthBin := strings.LastIndex( n, "strength" )
	intBin := strings.LastIndex( n, "interval" )
	mapped := make(map[string]string, 2*len(w.Attrs))
	if strengthBin >= 0 {
		for k, m := range w.Attrs {
				move := rand.Intn(len(m))
				mapped[k+"1"] = m[move]
			}
			for k, m := range w.Attrs {
				move := rand.Intn(len(m))
				mapped[k+"2"] = m[move]
			}
		} else if intBin >= 0 {
			t := w.Attrs["type"][rand.Intn(len(w.Attrs["type"]))]
			mapped["type"] = t
			var to_iter []string
			if t == "Repeat" {
				to_iter = []string{"time", "RepeatTime"}
			} else {
				to_iter = []string{"time", "LadderStartOrEnd", "LadderIntervals"}
				}
			for _, v := range to_iter {
					toinsert := w.Attrs[v][rand.Intn(len(w.Attrs[v]))]
					mapped[v] = toinsert
				}
		} else {
			for k, m := range w.Attrs {
				move := rand.Intn(len(m))
				mapped[k] = m[move]
			}
	}
	return Move{n, mapped}
}

func writeSpotifyPlaylist(w []Move) () {
	auth := spotify.NewAuthenticator("https://www.google.com", spotify.ScopeUserReadPrivate)
	url := auth.AuthURL("379780")
	fmt.Println(string(url))
}


func main() {
	text, err := readFile("Workouts.json")
	if err != nil {
		log.Fatalf("readFile: %s", err)
	}
	t := allWorkouts{}
	err = json.Unmarshal([]byte(text), &t)
	if err != nil {
		log.Fatalf("Unmarshal: %s", err)
	}
	skeleton := gen_workout_skeleton()
	var body []Move
	var to_append Move
	for _, move := range skeleton {
		for _, workout := range t.Workouts {
			if workout.Name == move {
				to_append = grab_random_workout(workout)
				body = append(body, to_append)
			}
		}
	}
	fmt.Println("JKR workout generator.")
	fmt.Println("Workout generated "+time.Now().String())
	for _, b := range body {
		fmt.Printf(b.Name+": ")
		fmt.Printf("[ \n")
		for k, v := range b.Attrs {
			fmt.Printf(k+": "+v+" \n")
	}
	fmt.Println("]")
	}
	writeSpotifyPlaylist(body)
}