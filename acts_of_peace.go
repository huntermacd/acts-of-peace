package actsOfPeace

import (
    "html/template"
    "net/http"

    "appengine"
    "appengine/datastore"
)

type ActOfPeace struct {
        Title string
        Description string
        FocusArea string
}

func init() {
    http.Handle("/static/", http.FileServer(http.Dir(".")))
    http.HandleFunc("/", root)
    http.HandleFunc("/submit", submit)
}

func actsOfPeaceKey(c appengine.Context) *datastore.Key {
        return datastore.NewKey(c, "Acts of Peace", "default_list_of_acts", 0, nil)
}

func root(w http.ResponseWriter, r *http.Request) {
        c := appengine.NewContext(r)
        q := datastore.NewQuery("ActOfPeace").Ancestor(actsOfPeaceKey(c))
        acts := make([]ActOfPeace, 0)
        if _, err := q.GetAll(c, &acts); err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
        }
        if err := actsOfPeaceTemplate.Execute(w, acts); err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
        }
}

var actsOfPeaceTemplate = template.Must(template.New("listOfActs").Parse(`
<html>
  <head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Submit Acts of Peace</title>
    <link rel="stylesheet" type="text/css" href="/static/style.css">
  </head>
  <body>
    <header>
      <img src="/static/logo.png" alt="Billion Acts logo.">
    </header>
    <h2 class="submit-header">Share Your Act of Peace</h2>
    <p>Fill out the form below telling us how you're helping to change the world.</p>
    <form action="/submit" method="post">
      <div>
        <label for="title">Title</label>
        <input type="text" name="title">
      </div>
      <div>
        <label for="description">Description</label>
        <textarea name="description" rows="3" cols="60"></textarea>
      </div>
      <div>
        <label for="focusArea">Focus Area</label>
        <select name="focusArea">
          <option value="Education & Community Development">Education & Community Development</option>
          <option value="Protecting the Environment">Protecting the Environment</option>
          <option value="Alleviating Extreme Poverty">Alleviating Extreme Poverty</option>
          <option value="Global Health & Wellness">Global Health & Wellness</option>
          <option value="Non-proliferation & Disarmament">Non-proliferation & Disarmament</option>
          <option value="Human Rights for All">Human Rights for All</option>
          <option value="Ending Racism & Hate">Ending Racism & Hate</option>
          <option value="Advancing Women & Children">Advancing Women & Children</option>
          <option value="Clean Water for Everyone">Clean Water for Everyone</option>
          <option value="Conflict Resolution">Conflict Resolution</option>
        </select>
      </div>
      <div>
        <input type="submit" value="Submit Act of Peace">
      </div>
    </form>
    <h2 class="submitted-header">See What Others Are Doing</h2>
    {{range .}}
      <div class="act">
        {{with .Title}}
          <h3>{{.}}</h3>
        {{end}}
        <p>{{.Description}}</p>
        <p><i>{{.FocusArea}}</i></p>
      </div>
    {{end}}
    <footer>
      PeaceJam <span class="year"></span>
    </footer>
    <script>
      var d = new Date;
      document.querySelector('.year').innerHTML = d.getFullYear();
    </script>
  </body>
</html>
`))

func submit(w http.ResponseWriter, r *http.Request) {
        c := appengine.NewContext(r)
        a := ActOfPeace{
                Title: r.FormValue("title"),
                Description: r.FormValue("description"),
                FocusArea: r.FormValue("focusArea"),
        }
        key := datastore.NewIncompleteKey(c, "ActOfPeace", actsOfPeaceKey(c))
        _, err := datastore.Put(c, key, &a)
        if err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
        }
        http.Redirect(w, r, "/", http.StatusFound)
}
