package course

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/OpenCoursePlatform/OpenCoursePlatform-Go/initiate"

	_ "github.com/go-sql-driver/mysql"
)

func TestGetSessionText(t *testing.T) {
	const SESSIONTEXT = "# Multis aprum sua propositique inde\n\n## Pro porrigis vocat acrior dis\n\nLorem markdownum perfudit harenam magna vetus illum reverentia, cum. Texere\nquodque, autem ambo ades ieiunia scopulum. Oris et visa mirabere paelice florem\no tempto **sensit** fidissime ut aegram videoque, adventu\n[viva](http://www.postquamquae.com/rumpitquoque) tulit pendet inplicuit,\ncalcitrat. Palamedes forte cura ter at trepidans tegit sinit: strenuitas.\n\n1. Est Pylio aut fratri acumina multiplicique pugnantem\n2. Soror saxo choro et vitae adstitit felix\n3. Ademptae longe sequentes ire\n\nQuamvis ait feres sequitur Iuppiter *adhuc rudibusque* est ast exit. Te tu\nrepentinos functus ululatibus agunt, at Nilus lacerata flammasque: interea.\nAltissima Emathion tamen rigorem et naribus: et avertere, priscis **victricemque\nCereale** aequorei longisque dummodo lapsis gratia, sit.\n\n> *Sed sanguine adiere* dominum [pietate](http://cicatrixaudire.io/) restat,\n> sole nigram, dimittit. Coniunx limbus ad oppida si crinem auctor clamore modo\n> inertes. **Et resisti** devicto verso, erat tela, die exosus, alios imo\n> quondam non Aiaci. De tutum nomine spectabat est dictis dona, vestigia\n> lanugine poscat octavo loton eo optas: tenens.\n\n## Ipse disponunt regia obliqua patet\n\nSanguine senserat fontesque addit figuras tulerunt, male fatigat me dura\nnavifragumque lauro adspicio. Cervice illos latissima: sed rotat humana fore;\nflumina.\n\nNec haec maxima? Possideat exuit gratissime Troas et placandam voces Olympum\ntoris.\n\n- Nisi posita summis\n- Figuram horruit Lucina\n- Iove potest urbe facientia declivis\n- Dixit licet suis medio at camini\n- Ante voce custodia magnae\n- Etiamnunc visceribus visis admissumque illo onere me\n\nNatantibus et clamore ferroque, de opus. Ignarum diu, et mea reliquit vocant.\n*Et doloris in* nido!\n\nAdest et, haec, unda aer ululasse [Venus](http://senexte.com/verbis). Falleret\naeno mandata, qui desubito, verum captaeque illo imagine pharetraque enim\nnequiquam annis?"
	db, err := initiate.Tests()
	if err != nil {
		t.Errorf("Initiation of tests failed. Error message %s", err.Error())
		return
	}

	_, err = db.Exec(`
	INSERT INTO session (id, name, slug, module_id, session_type, created_at)
	VALUES
		(1, 'Introduction!', 'introduction', 1, 0, '2019-11-13 22:12:55');
	
	`)
	if err != nil {
		t.Errorf("Insertion of session in database failed. Error message: %s", err.Error())
		return
	}

	_, err = db.Exec(`
	INSERT INTO session_text (id, text, session_id, created)
	VALUES
		(1, ?, 1, '2019-11-14 13:35:27');	
	`, SESSIONTEXT)
	if err != nil {
		t.Errorf("Insertion of session in database failed. Error message: %s", err.Error())
		return
	}

	sessionText, err := GetSessionsText(db, 1)
	if err != nil {
		t.Errorf("GetSessions failed. Error message: %s", err.Error())
		return
	}

	if sessionText != SESSIONTEXT {
		t.Errorf("GetSessions failed. Texts are not equal")
		return
	}

	if sessionText == "SESSIONTEXT" {
		t.Errorf("GetSessions failed. Texts should not be equal")
		return
	}

	err = initiate.FinishTests(db)
	if err != nil {
		t.Errorf("Finishing of tests failed. Error message %s", err.Error())
		return
	}
}

func TestGetSession(t *testing.T) {
	db, err := initiate.Tests()
	if err != nil {
		t.Errorf("Initiation of tests failed. Error message %s", err.Error())
		return
	}

	_, err = db.Exec(`
	INSERT INTO session (id, name, slug, module_id, session_type, created_at)
	VALUES
		(1, 'Introduction!', 'introduction', 1, 1, '2019-11-13 22:12:55');
	
	`)
	if err != nil {
		t.Errorf("Insertion of session in database failed. Error message: %s", err.Error())
		return
	}

	session, moduleID, err := GetSession(db, "introduction")
	if err != nil {
		t.Errorf("GetSession failed. Error message: %s", err.Error())
		return
	}

	if moduleID != 1 {
		t.Errorf("GetSession failed. ModuleID are not equal")
		return
	}

	sessionCopy := Session{ID: 1, Name: "Introduction!", SessionType: 1}

	if session != sessionCopy {
		t.Errorf("GetSession failed. Sessions are not equal")
		return
	}

	err = initiate.FinishTests(db)
	if err != nil {
		t.Errorf("Finishing of tests failed. Error message %s", err.Error())
		return
	}
}

func TestGetSessionsYoutube(t *testing.T) {
	const SESSIONTEXT = "# Multis aprum sua propositique ind\r\n\r\n## Pro porrigis vocat acrior di\r\n"
	db, err := initiate.Tests()
	if err != nil {
		t.Errorf("Initiation of tests failed. Error message %s", err.Error())
		return
	}

	_, err = db.Exec(`
	INSERT INTO session_youtube (id, text, youtube_id, session_id, created)
	VALUES
		(1, '# Multis aprum sua propositique ind\r\n\r\n## Pro porrigis vocat acrior di\r\n', 'dBn6_qflUZ4', 3, '2019-11-14 13:35:27');	
	`)
	if err != nil {
		t.Errorf("Insertion of session in database failed. Error message: %s", err.Error())
		return
	}

	sessionText, youtubeID, err := GetSessionsYoutube(db, 3)
	if err != nil {
		t.Errorf("GetSession failed. Error message: %s", err.Error())
		return
	}
	if sessionText != SESSIONTEXT {
		t.Errorf("GetSession failed. Texts are not equal.")
		return
	}

	if youtubeID != "dBn6_qflUZ4" {
		t.Errorf("GetSession failed. Youtube IDs are not equal.")
		return
	}

	err = initiate.FinishTests(db)
	if err != nil {
		t.Errorf("Finishing of tests failed. Error message %s", err.Error())
		return
	}
}

func TestSessionPage(t *testing.T) {
	req, err := http.NewRequest("GET", "", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(SessionPage)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}
}

func TestDBForSessionPage(t *testing.T) {
	config, err := initiate.DeleteSettingsFile()
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(SessionPage)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}

	expected := ""
	if rr.Body.String() == expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}

	err = initiate.WriteSettingsFile(config)
	if err != nil {
		t.Fatal(err)
	}

}
