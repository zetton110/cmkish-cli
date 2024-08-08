package action

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
	model "github.com/zetton110/cmkish-cli/model"
)

type FindSong struct {
	Title        string
	ProgramTitle string
	Artist       string
	DatabasePath string
}

func (f *FindSong) Run() error {
	title := f.Title
	programTitle := f.ProgramTitle
	artist := f.Artist
	databasePath := f.DatabasePath

	db, err := sql.Open("sqlite3", databasePath)
	if err != nil {
		return err
	}
	defer db.Close()

	var queries []string
	if len(programTitle) > 0 {
		queries = append(queries, fmt.Sprintf("SELECT title, artist, program_name, op_ed, broadcast_order FROM anison INNER JOIN program ON anison.program_id = program.ID where program.name LIKE '%%%s%%' ORDER BY program.start_date ASC", programTitle))
		queries = append(queries, fmt.Sprintf("SELECT title, artist, program_name, op_ed, broadcast_order FROM game INNER JOIN program ON game.program_id = program.ID where program.name LIKE '%%%s%%' ORDER BY program.start_date ASC", programTitle))
		queries = append(queries, fmt.Sprintf("SELECT title, artist, program_name, op_ed, broadcast_order FROM side_effect INNER JOIN program ON side_effect.program_id = program.ID where program.name LIKE '%%%s%%' ORDER BY program.start_date ASC", programTitle))
	} else if len(title) > 0 {
		queries = append(queries, fmt.Sprintf("SELECT title, artist, program_name, op_ed, broadcast_order FROM anison where title LIKE '%%%s%%'", title))
		queries = append(queries, fmt.Sprintf("SELECT title, artist, program_name, op_ed, broadcast_order FROM game where title LIKE '%%%s%%'", title))
		queries = append(queries, fmt.Sprintf("SELECT title, artist, program_name, op_ed, broadcast_order FROM side_effect where title LIKE '%%%s%%'", title))
	} else if len(artist) > 0 {
		queries = append(queries, fmt.Sprintf("SELECT title, artist, program_name, op_ed, broadcast_order FROM anison where artist LIKE '%%%s%%'", artist))
		queries = append(queries, fmt.Sprintf("SELECT title, artist, program_name, op_ed, broadcast_order FROM game where artist LIKE '%%%s%%'", artist))
		queries = append(queries, fmt.Sprintf("SELECT title, artist, program_name, op_ed, broadcast_order FROM side_effect where artist LIKE '%%%s%%'", artist))
	}

	var songs []model.Song
	for _, q := range queries {

		rows, err := db.Query(q)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var s model.Song
			err := rows.Scan(&s.Title, &s.Artist, &s.ProgramName, &s.OpEd, &s.BroadcastOrder)
			if err != nil {
				fmt.Errorf("failed to parse anison. %w\n", err)
			}
			songs = append(songs, s)
		}
	}

	if len(songs) == 0 {
		fmt.Println("Nothig is found.")
		return nil
	}

	data := [][]string{}
	for _, a := range songs {
		data = append(data, []string{
			a.Title,
			a.Artist,
			a.ProgramName,
			a.OpEd + " " + a.BroadcastOrder,
		})
	}
	header := []string{"曲名", "歌手", "作品名", "備考"}

	renderTable(data, header)

	fmt.Printf("%d hits.\n", len(songs))

	return nil
}

func renderTable(data [][]string, header []string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(header)

	for _, v := range data {
		table.Append(v)
	}
	table.Render()

}