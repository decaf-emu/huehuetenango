package ghs

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

var templatePrefixes = []string{"tm", "ps", "pt"}

var baseNames = map[string]string{
	"__vtbl": " virtual table",
	"__ct":   "#",
	"__dt":   "~#",
	"__as":   "operator=",
	"__eq":   "operator==",
	"__ne":   "operator!=",
	"__gt":   "operator>",
	"__lt":   "operator<",
	"__ge":   "operator>=",
	"__le":   "operator<=",
	"__pp":   "operator++",
	"__pl":   "operator+",
	"__apl":  "operator+=",
	"__mi":   "operator-",
	"__ami":  "operator-=",
	"__ml":   "operator*",
	"__amu":  "operator*=",
	"__dv":   "operator/",
	/* XXX below baseNames have not been seen - guess from libiberty cplus-dem.c*/
	"__adv": "operator/=",
	"__nw":  "operator new",
	"__dl":  "operator delete",
	"__vn":  "operator new[]",
	"__vd":  "operator delete[]",
	"__md":  "operator%",
	"__amd": "operator%=",
	"__mm":  "operator--",
	"__aa":  "operator&&",
	"__oo":  "operator||",
	"__or":  "operator|",
	"__aor": "operator|=",
	"__er":  "operator^",
	"__aer": "operator^=",
	"__ad":  "operator&",
	"__aad": "operator&=",
	"__co":  "operator~",
	"__cl":  "operator",
	"__ls":  "operator<<",
	"__als": "operator<<=",
	"__rs":  "operator>>",
	"__ars": "operator>>=",
	"__rf":  "operator->",
	"__vc":  "operator[]",
}

var baseTypes = map[byte]string{
	'v': "void",
	'i': "int",
	's': "short",
	'c': "char",
	'w': "wchar_t",
	'b': "bool",
	'f': "float",
	'd': "double",
	'l': "long",
	'L': "long long",
	'e': "...",
	/* XXX below baseTypes have not been seen - guess from libiberty cplus-dem.c */
	'r': "long double",
}

var typePrefixes = map[byte]string{
	'U': "unsigned",
	'S': "signed",
	/* XXX below typePrefixes have not been seen - guess from libiberty cplus-dem.c */
	'J': "__complex",
}

var typeSuffixes = map[byte]string{
	'P': "*",
	'R': "&",
	'C': "const",
	'V': "volatile", /* XXX this is a guess! */
	/* XXX below typeSuffixes have not been seen - guess from libiberty cplus-dem.c */
	'u': "restrict",
}

func Demangle(input string) (string, error) {
	name, err := decompress(input)
	if err != nil {
		return input, err
	}

	baseName, mangle, err := readBaseName(name)
	if err != nil {
		return input, err
	}

	var declStatic = ""
	if strings.HasPrefix(mangle, "S__") {
		declStatic = "static "
		mangle = mangle[3:]
	}

	var declNameSpace string
	var declClass string

	if strings.HasPrefix(mangle, "Q") {
		declNameSpace, mangle, err = readNameSpace(mangle)
		if err != nil {
			return input, err
		}

		var last = strings.LastIndex(declNameSpace, "::")
		if last == -1 {
			declClass = declNameSpace
		} else {
			declClass = declNameSpace[last+2:]
		}

		declNameSpace += "::"
	} else if len(mangle) > 0 && startsWithDigit(mangle) {
		declClass, mangle, err = readString(mangle)
		if err != nil {
			return input, err
		}
		declNameSpace = declClass + "::"
	}

	baseName = strings.Replace(baseName, "#", declClass, -1)

	if strings.HasPrefix(mangle, "S") {
		declStatic = "static "
		mangle = mangle[1:]
	}

	var declConst string
	if strings.HasPrefix(mangle, "C") {
		declConst = " const"
		mangle = mangle[1:]
	}

	declType := "#"
	if strings.HasPrefix(mangle, "F") {
		declType, mangle, err = readType(nil, mangle)
		if err != nil {
			return input, err
		}
	}

	if strings.HasPrefix(mangle, "_") {
		var end, err = strconv.Atoi(mangle[1:])
		if err != nil {
			return input, err
		}
		baseName += "_" + strconv.Itoa(end)
		mangle = ""
	}

	if len(mangle) > 0 {
		return input, fmt.Errorf("Unknown modifier %v", mangle)
	}

	declType = strings.Replace(declType, "(#)", " "+declNameSpace+baseName, -1)
	declType = strings.Replace(declType, "#", declNameSpace+baseName, -1)
	combination := declStatic + declType + declConst

	return strings.Replace(combination, "::"+baseNames["__vtbl"], baseNames["__vtbl"], -1), nil
}

func startsWithDigit(input string) bool {
	if len(input) == 0 {
		return false
	}
	r, size := utf8.DecodeRuneInString(input[0:1])
	if size > 0 {
		return unicode.IsDigit(r)
	}
	return false
}

func decompress(name string) (string, error) {
	if !strings.HasPrefix(name, "__CPR") {
		return name, nil
	}

	name = name[5:] //skip '__CPR'

	decompressedLen, name, err := readIntPrefix(name)
	if err != nil {
		return "", err
	}

	name = name[2:] //skip '__'
	segments := strings.Split(name, "J")

	var result = ""
	for i, val := range segments {
		//I assume, perhaps wrongly, that even elements are literals
		if i%2 == 0 {
			//Literal
			result += val
		} else {
			//Interpolation
			var offset = -1
			offset, err = strconv.Atoi(val)
			if offset == -1 || err != nil {
				return "", fmt.Errorf("Bad decompression offset. %v is not a valid offset.", val)
			}

			var sub, _, err = extractName(result[offset:])
			if err != nil {
				return "", err
			}
			result += writeString(sub)
		}
	}

	if len(result) != decompressedLen {
		return "", fmt.Errorf("Bad decompression length. Expected %v and got %v.", decompressedLen, len(result))
	}

	return result, nil
}

func demangleTemplate(name string) (string, error) {
	mstart := strings.Index(name[1:], "__")
	if mstart != -1 && strings.HasPrefix(name[mstart:], "___") {
		mstart++
	}

	if mstart == -1 {
		return name, nil
	}

	remainder := name[mstart+3:]
	name = name[:mstart+1]

	for {
		var err error
		if !startsWithAny(remainder, templatePrefixes) {
			return "", errors.New("Unexpected template argument prefix.")
		}

		var lstart = strings.Index(remainder, "__")
		if lstart == -1 {
			return "", errors.New("Bad template argument")
		}

		remainder = remainder[lstart+2:]
		var extracted = ""
		extracted, remainder, err = extractName(remainder)
		if err != nil {
			return "", errors.New("Bad template argument length.")
		}

		if !strings.HasPrefix(extracted, "_") {
			return "", fmt.Errorf("Unexpected character after template parameter length. \"%v\"", extracted)
		}

		var declArgs = ""
		var tmp = ""
		declArgs, tmp, err = readTemplateArguments(extracted[1:])
		if err != nil {
			return "", err
		}

		if strings.HasSuffix(declArgs, ">") {
			declArgs += " "
		}

		name += "<" + declArgs + ">"

		if tmp != remainder {
			return "", fmt.Errorf("Bad template argument length.  %s != %s", tmp, remainder)
		}

		if len(remainder) == 0 {
			return name, nil
		}

		if !strings.HasPrefix(remainder, "__") {
			return "", fmt.Errorf("Unexpected character(s) after template: %c%c.  Expected \"__\".", remainder[0], remainder[1])
		}

		remainder = remainder[2:]
	}
	return "", errors.New("Unexpectedly exited outside unbounded loop")
}

func readBaseName(name string) (string, string, error) {
	if len(name) == 0 {
		return "", "", errors.New("Unexpected end of string, Expected a name.")
	}

	var opName = ""
	if strings.HasPrefix(name, "__op") {
		var t, name, err = readType(nil, name)
		if err != nil {
			return "", "", err
		}

		opName = "operator " + t
		name = "#" + name
	}

	var mstart = strings.Index(name[1:], "__")
	if mstart != -1 && strings.HasPrefix(name[mstart:], "___") {
		mstart++
	}

	if mstart == -1 {
		return name, "", nil
	}

	var remainder = ""
	name, remainder = name[:mstart+1], name[mstart+3:]

	if val, ok := baseNames[name]; ok {
		name = val
	} else if name == "#" {
		name = opName
	}

	for startsWithAny(remainder, templatePrefixes) {
		var lstart = strings.Index(remainder, "__")
		if lstart == -1 {
			return "", "", errors.New("Bad template argument")
		}

		name += "__" + remainder[:lstart]
		remainder = remainder[lstart+2:]

		var name, remainder, err = extractName(remainder)
		if err != nil {
			return "", "", err
		}

		name += "__" + writeString(name)
		if len(remainder) == 0 {
			return name, remainder, nil
		}

		if !strings.HasPrefix(remainder, "__") {
			return "", "", errors.New("Unexpected character(s) after template.")
		}

		remainder = remainder[2:]
	}

	var dt, err = demangleTemplate(name)
	if err != nil {
		return "", "", err
	}

	return dt, remainder, nil
}

func readArguments(name string) (string, string, error) {
	var result = ""
	var args = []string{}
	var remainder = name
	var err error = nil

	for len(remainder) > 0 && !strings.HasPrefix(remainder, "_") {
		if len(args) > 0 {
			result += ", "
		}

		var t = ""
		t, remainder, err = readType(args, remainder)
		if err != nil {
			return "", "", err
		}

		result += strings.Replace(t, "#", "", -1)
		args = append(args, t)
	}

	return result, remainder, nil
}

func readTemplateArguments(input string) (string, string, error) {
	var result = ""
	var args = []string{}
	var remainder = input

	var tipe = "" //Can't call it 'type' as the original source does :)
	var val = ""
	for len(remainder) > 0 && !strings.HasPrefix(remainder, "_") {
		var err error = nil
		if len(args) > 0 {
			result += ", "
		}

		if strings.HasPrefix(remainder, "X") {
			remainder = remainder[1:]
			if len(remainder) == 0 {
				return "", "", errors.New("Unexpected end of string.  Expected a type.")
			}

			if startsWithDigit(remainder) {
				tipe = "#"
				val, remainder, err = readString(remainder)
				if err != nil {
					return "", "", err
				}
			} else {
				tipe, remainder, err = readType(args, remainder)
				if err != nil {
					return "", "", err
				}
				tipe = strings.Replace(tipe, "#", " #", -1)

				if strings.HasPrefix(remainder, "L") {
					remainder = remainder[1:] //skip 'L'
					if len(remainder) == 0 {
						return "", "", errors.New("Unexpected end of string.  Expected \"_\".")
					}
					if !strings.HasPrefix(remainder, "_") {
						return "", "", fmt.Errorf("Unexpected character after template parameter encoding %c.  Expected \"_\".", remainder[0])
					}

					var length = 0
					length, remainder, err = readIntPrefix(remainder[1:])
					if err != nil {
						return "", "", errors.New("Bad template parameter length.")
					}

					if !strings.HasPrefix(remainder, "_") {
						return "", "", fmt.Errorf("Unexpected character after template parameter length. \"%v\"", remainder)
					}

					remainder = remainder[1:] //skip '_'
					val = remainder[:length]
					remainder = remainder[length:]
				} else {
					return "", "", fmt.Errorf("Unknown template parameter encoding \"%c\".", remainder[0])
				}
			}
		} else {
			val, remainder, err = readType(args, remainder)
			if err != nil {
				return "", "", err
			}
			tipe = "class #"
		}

		result += strings.Replace(tipe, "#", "Z"+strconv.Itoa(len(args)+1)+" = "+val, -1)
		args = append(args, val)
	}

	return result, remainder, nil
}

func readType(args []string, name string) (string, string, error) {
	if len(name) == 0 {
		return "", "", errors.New("Unexpected end of string.  Expected a type.")
	}

	if val, ok := baseTypes[name[0]]; ok {
		return val + "#", name[1:], nil
	} else if strings.HasPrefix(name, "Q") {
		var result, remainder, err = readNameSpace(name)
		if err != nil {
			return "", "", err
		}
		return result + "#", remainder, nil
	} else if startsWithDigit(name) {
		var result, remainder, err = readString(name)
		if err != nil {
			return "", "", err
		}
		return result + "#", remainder, nil
	} else if val, ok := typePrefixes[name[0]]; ok {
		var result, remainder, err = readType(args, name[1:])
		if err != nil {
			return "", "", err
		}
		return val + " " + result, remainder, nil
	} else if val, ok := typeSuffixes[name[0]]; ok {
		var result, remainder, err = readType(args, name[1:])
		if err != nil {
			return "", "", err
		}
		return strings.Replace(result, "#", " "+val+"#", -1), remainder, nil
	} else if strings.HasPrefix(name, "Z") {
		var index = strings.Index(name[1:], "Z") //next 'Z'
		if index == -1 {
			return "", "", errors.New("Unexpected end of string. Expected \"Z\"")
		}
		return name[:index] + "#", name[:index+1], nil
	} else if strings.HasPrefix(name, "A") {
		name = name[1:]
		/*
			var index = strings.Index(name, "_Z")
			if index > -1 {
				var end = strings.Index(name[2:], "Z") //next 'Z'
				length = strconv.Atoi(name[1 : end-1])
				name = name[end+1:]
			} else {
		*/
		var length, remainder, err = readIntPrefix(name)

		if err != nil || len(remainder) == 0 {
			return "", "", fmt.Errorf("Unexpected end of string.  Expected \"_\".")
		}
		if !strings.HasPrefix(remainder, "_") {
			return "", "", fmt.Errorf("Unexpected character after array length \"%c\".  Expected \"_\".", remainder[0])
		}
		remainder = remainder[1:] //skip over '_'
		var result = ""
		result, remainder, err = readType(args, remainder)
		if err != nil {
			return "", "", err
		}

		return strings.Replace(result, "#", "#["+strconv.Itoa(length)+"]", -1), remainder, nil
	} else if strings.HasPrefix(name, "F") {
		var declArgs, name, err = readArguments(name[1:])
		if err != nil {
			return "", "", err
		}
		if args == nil && (len(name) == 0 || strings.HasPrefix(name, "_")) {
			return "#(" + declArgs + ")", name, nil
		}
		if len(name) == 0 {
			return "", "", errors.New("Unexpected end of string, expected \"_\".")
		}
		if !strings.HasPrefix(name, "_") {
			return "", "", fmt.Errorf("Unexpected character after template parameter length. \"%v\"", name)
		}
		var result, remainder = "", ""
		result, remainder, err = readType(args, name[1:])
		if err != nil {
			return "", "", err
		}

		return strings.Replace(result, "#", "(#)("+declArgs+")", -1), remainder, nil
	} else if strings.HasPrefix(name, "T") {
		var index, err = strconv.Atoi(name[1:2])
		if err != nil {
			return "", "", err
		}
		if len(args) < index {
			return "", "", fmt.Errorf("Bad argument number \"%v\".", index)
		}

		return args[index-1], name[2:], nil
	} else if strings.HasPrefix(name, "N") {
		var err error = nil
		var count, arg = 0, 0
		var remainder = ""
		count, err = strconv.Atoi(name[1:2])
		if err != nil {
			return "", "", err
		}
		arg, err = strconv.Atoi(name[2:3])
		if err != nil {
			return "", "", err
		}

		if count > 1 {
			remainder = "N" + strconv.Itoa(count-1) + strconv.Itoa(arg) + name[3:]
		} else {
			remainder = name[3:]
		}

		return args[arg-1], remainder, nil
	}

	return "", "", fmt.Errorf("Unknown type \"%c\".", name[0])
}

func readNameSpace(input string) (string, string, error) {
	if len(input) == 0 || !strings.HasPrefix(input, "Q") {
		return "", "", errors.New("Unexpected end of string.  Expected \"Q\".")
	}

	//Q2_5Types58PortionedHandle__tm__35_5AgentXCUiL_2_14XCUiL_1_6XCUiL_1_0b
	var namespaces = []string{}
	var count, remainder, err = readIntPrefix(input[1:])
	if err != nil {
		return "", "", err
	}

	remainder = remainder[1:] //step over '_'

	for i := 0; i < count; i++ {
		//var remainderspaces = strings.SplitAfter(remainder, "Z")
		var ns = ""
		ns, remainder, err = readString(remainder)
		if err != nil {
			return "", "", err
		}
		namespaces = append(namespaces, ns)
	}

	return strings.Join(namespaces, "::"), remainder, nil
}

func readString(input string) (string, string, error) {
	if len(input) == 0 {
		return "", "", errors.New("Unexpected end of string.  Expected a digit.")
	}

	var name, remainder, err = extractName(input)
	if err != nil {
		return "", "", err
	}

	var dt = ""
	dt, err = demangleTemplate(name)
	if err != nil {
		return "", "", err
	}

	return dt, remainder, nil
}

func writeString(input string) string {
	//Prefix string with its length.
	return strconv.Itoa(len(input)) + input
}

func readIntPrefix(input string) (int, string, error) {
	if len(input) == 0 {
		return -1, "", errors.New("Unexpected end of string.  Expected a digit.")
	}

	re := regexp.MustCompile(`([0-9]+)(.*)$`)
	var results = re.FindStringSubmatch(input)

	//Should have 3 items in array; the whole regex match, and then each capture.
	if results == nil || len(results) < 3 {
		return -1, "", errors.New("Unexpected end of string.  Unable to match digits.")
	}

	var length, err = strconv.Atoi(results[1])
	if err != nil {
		return 0, "", err
	}
	var postNumber = results[2]
	return length, postNumber, nil
}

func extractName(name string) (string, string, error) {
	length, postNumber, err := readIntPrefix(name)
	return postNumber[:length], postNumber[length:], err
}

func startsWithAny(input string, names []string) bool {
	for _, v := range names {
		if strings.HasPrefix(input, v) {
			return true
		}
	}
	return false
}
