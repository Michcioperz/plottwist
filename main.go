package main

import (
	"fmt"
	"github.com/rthornton128/goncurses"
	"log"
	"os/exec"
)

func submenu(scr *goncurses.Window, menu *goncurses.Menu) {
	curr := menu.Current(nil)
	episodesPage, err := FetchPageContents(UrlPseudoJoin("/a/" + curr.Name()))
	if err != nil {
		scr.Print(err)
		return
	}
	episodes, err := ExtractEpisodesList(episodesPage)
	if err != nil {
		scr.Print(err)
		return
	}
	deepMenuItems := make([]*goncurses.MenuItem, len(episodes)+1)
	for i, episode := range episodes {
		deepMenuItems[i], _ = goncurses.NewItem(fmt.Sprintf("Episode %v", episode.Number), UrlPseudoJoin(episode.Source))
		defer deepMenuItems[i].Free()
	}
	deepMenuItems[len(episodes)], _ = goncurses.NewItem("Play all", UrlPseudoJoin("/a/"+curr.Name()))
	defer deepMenuItems[len(episodes)].Free()
	defer scr.Refresh()
	menu.UnPost()
	defer menu.Post()
	depth := true
	deepMenu, err := goncurses.NewMenu(deepMenuItems)
	if err != nil {
		scr.Print(err)
		return
	}
	defer deepMenu.Free()
	deepBuffer := []rune("")
	deepMenu.SetWindow(scr)
	deepMenu.Post()
	scr.Refresh()
	for depth {
		goncurses.Update()
		ch := scr.GetChar()
		chch := goncurses.KeyString(ch)
		switch chch {
		case "tab":
			deepMenu.Driver(goncurses.REQ_NEXT_MATCH)
		case "backspace":
			if len(deepBuffer) > 0 {
				deepBuffer = deepBuffer[:len(deepBuffer)-1]
				deepMenu.SetPattern(string(deepBuffer))
			} else {
				depth = false
				continue
			}
		case "down":
			deepMenu.Driver(goncurses.REQ_DOWN)
		case "up":
			deepMenu.Driver(goncurses.REQ_UP)
		case "page down":
			deepMenu.Driver(goncurses.REQ_PAGE_DOWN)
		case "page up":
			deepMenu.Driver(goncurses.REQ_PAGE_UP)
		case "home":
			deepMenu.Driver(goncurses.REQ_FIRST)
		case "end":
			deepMenu.Driver(goncurses.REQ_LAST)
		case "enter":
			ep := deepMenu.Current(nil)
			ep.SetValue(true)
			goncurses.Update()
			exec.Command("mpv", "--fs", ep.Description()).Run()
			ep.SetValue(false)
		default:
			deepBuffer = append(deepBuffer, rune(ch))
			deepMenu.SetPattern(string(deepBuffer))
		}
	}
}

func main() {
	scr, err := goncurses.Init()
	if err != nil {
		log.Panic("error on ncurses init", err)
	}
	goncurses.Raw(true)
	goncurses.Echo(false)
	goncurses.Cursor(0)
	scr.Clear()
	scr.Keypad(true)
	scr.Print("any letter typed appends to search query\ntab goes to next matched item\nbackspace deletes the last letter or quits when query is empty\nenter, pgup, pgdown, home, end and arrows work as usual")
	goncurses.Update()
	scr.Refresh()
	defer goncurses.End()
	listHtml, err := FetchPageContents(TwistRoot)
	if err != nil {
		scr.Print(err)
		return
	}
	series, err := ExtractSeriesList(listHtml)
	if err != nil {
		scr.Print(err)
		return
	}
	menuItems := make([]*goncurses.MenuItem, len(series))
	for i, serie := range series {
		menuItems[i], _ = goncurses.NewItem(serie.Slug, serie.NiceTitle())
		defer menuItems[i].Free()
	}
	menu, err := goncurses.NewMenu(menuItems)
	if err != nil {
		scr.Print(err)
		return
	}
	defer menu.Free()
	menu.SetWindow(scr)
	menu.Post()
	buffer := []rune("")
	for {
		goncurses.Update()
		ch := scr.GetChar()
		chch := goncurses.KeyString(ch)
		switch chch {
		case "tab":
			menu.Driver(goncurses.REQ_NEXT_MATCH)
		case "backspace":
			if len(buffer) > 0 {
				buffer = buffer[:len(buffer)-1]
				menu.SetPattern(string(buffer))
			} else {
				return
			}
		case "down":
			menu.Driver(goncurses.REQ_DOWN)
		case "up":
			menu.Driver(goncurses.REQ_UP)
		case "page down":
			menu.Driver(goncurses.REQ_PAGE_DOWN)
		case "page up":
			menu.Driver(goncurses.REQ_PAGE_UP)
		case "enter":
			submenu(scr, menu)
		case "home":
			menu.Driver(goncurses.REQ_FIRST)
		case "end":
			menu.Driver(goncurses.REQ_LAST)
		default:
			buffer = append(buffer, rune(ch))
			menu.SetPattern(string(buffer))
		}
	}
}
