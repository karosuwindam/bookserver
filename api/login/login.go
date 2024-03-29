package login

import (
	"bookserver/api/common"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"time"
	"unsafe"

	"github.com/dgrijalva/jwt-go"
)

// UName() ユーザの名前
func (t *JwtData) UName() string { return t.userName }

// UType() ユーザーの種類
func (t *JwtData) UType() UserType { return t.usertype }

// Ext() 時間出力
func (t *JwtData) Ext() time.Time { return time.Unix(t.ext, 0) }

func cklogin(name, pass string) bool {
	if psckdata[name] == pass && pass != "" {
		return true
	}
	return false
}

func convetStructoToMap(str interface{}) (jwt.MapClaims, error) {
	out := jwt.MapClaims{}
	sv := reflect.ValueOf(str).Elem()
	if reflect.ValueOf(str).Kind() != reflect.Ptr {
		return out, errors.New("input data is not pointer")
	}
	// ssv := reflect.ValueOf(sv.Interface())
	st := reflect.TypeOf(sv.Interface())
	for i := 0; i < st.NumField(); i++ {
		f := st.Field(i)
		fv := sv.FieldByName(f.Name)
		if tag := f.Tag.Get("json"); tag != "" {
			fv = reflect.NewAt(fv.Type(), unsafe.Pointer(fv.UnsafeAddr())).Elem()

			v := fv.Interface()
			switch v.(type) {
			case string:
				out[tag] = interface{}(v.(string))
			case int64:
				out[tag] = interface{}(v.(int64))
			case UserType:
				out[tag] = interface{}(int64(v.(UserType)))
			}
		}
	}
	return out, nil
}

func createjwt(name, pass string) (string, error) {
	if cklogin(name, pass) {
		claims, err := convetStructoToMap(&JwtData{userName: name, ext: time.Now().Add(time.Hour * 24).Unix(), usertype: ADMIN | GUEST | USER})
		if err != nil {

		}
		// ヘッダーとペイロードの生成
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString([]byte(jwtsecretkey))
		if err != nil {
			return "", err
		}
		return tokenString, nil

	}

	return "", nil
}

func convetMapToStruct(m map[string]interface{}, str interface{}) error {
	if reflect.ValueOf(str).Type().Kind() != reflect.Ptr {
		return errors.New("Don't struct pointer input i=" + reflect.ValueOf(str).Type().Kind().String())
	}

	if len(m) == 0 {
		return nil
	}
	sv := reflect.ValueOf(str).Elem()
	st := reflect.TypeOf(str).Elem()

	for i := 0; i < st.NumField(); i++ {
		ft := st.Field(i)
		v := sv.FieldByName(ft.Name)
		v = reflect.NewAt(v.Type(), unsafe.Pointer(v.UnsafeAddr())).Elem()
		s := m[ft.Tag.Get("json")]
		switch ft.Type.Kind() {
		case reflect.Int64:
			v.SetInt(int64(s.(float64)))
		case reflect.String:
			v.SetString(s.(string))
		case reflect.Int:
			v.SetInt(int64(s.(float64)))
		}
	}
	return nil

}

func Unpackjwt(jwtdata string) (JwtData, error) {
	var output JwtData
	token, err := jwt.Parse(jwtdata, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(jwtsecretkey), nil
	})
	if err != nil {
		return output, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if err := convetMapToStruct(claims, &output); err != nil {
			return output, err
		}

	}
	return output, nil

}

func login(w http.ResponseWriter, r *http.Request) {
	ses, _ := cs.Get(r, "hello-session")
	ses.Values["login"] = nil
	ses.Values["name"] = nil
	ses.Values["time"] = nil
	nm := r.PostFormValue("name")
	pw := r.PostFormValue("pass")
	if cklogin(nm, pw) {
		ses.Values["login"] = true
		ses.Values["name"] = nm
		ses.Values["time"] = time.Now().Add(time.Hour).Unix()
		if jwtdata, err := createjwt(nm, pw); err == nil {
			jwttmp[nm] = jwtdata
		} else {
			jwttmp[nm] = ""
		}
	}
	ses.Save(r, w)
	getlogin(w, r)
}

func TokenReNew(w http.ResponseWriter, r *http.Request) {
	ses, _ := cs.Get(r, "hello-session")
	name, _ := ses.Values["name"].(string)
	ses.Values["time"] = time.Now().Add(time.Hour).Unix()
	claims, err := convetStructoToMap(&JwtData{userName: name, ext: time.Now().Add(time.Hour * 24).Unix(), usertype: ADMIN | GUEST | USER})
	if err != nil {

	}
	// ヘッダーとペイロードの生成
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte(jwtsecretkey))
	jwttmp[name] = tokenString
	ses.Save(r, w)

}

func getlogin(w http.ResponseWriter, r *http.Request) {
	msg := common.Result{Name: "login", Code: http.StatusOK, Date: time.Now()}
	ses, _ := cs.Get(r, "hello-session")
	flg, _ := ses.Values["login"].(bool)
	nm, _ := ses.Values["name"].(string)
	limittime, _ := ses.Values["time"].(int64)
	if flg {
		if token := jwttmp[nm]; token != "" && limittime >= time.Now().Unix() {
			TokenReNew(w, r)
			msg.Result = tResurlt{status: "User:" + nm + " " + "LOGIN OK", Token: token}
		} else {
			msg.Result = "User:" + nm + " " + "LOGIN OK"

		}
	} else {
		msg.Code = http.StatusUnauthorized
		msg.Result = "User:" + nm + " " + "LOGIN NG"
	}
	fmt.Println(r.Method, r.URL.Path, nm, flg)
	common.CommonBack(msg, w)

}

func CkUserlogin(name string, r *http.Request) UserType {
	ses, _ := cs.Get(r, "hello-session")
	flg, _ := ses.Values["login"].(bool)
	if flg {
		return ADMIN | GUEST | USER
	}
	return GUEST
}

func webServerLogin(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		login(w, r)
	default:
		getlogin(w, r)
	}
}
