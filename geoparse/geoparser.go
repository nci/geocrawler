package main

import (
	"../geolib"
	//geo "bitbucket.org/monkeyforecaster/geometry"
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"time"
)

type GeoMetaData struct {
	DataSetName    string            `json:"ds_name"`
	TimeStamps     []time.Time       `json:"timestamps"`
	FileNameFields map[string]string `json:"filename_fields"`
	Polygon        string            `json:"polygon"`
	RasterCount    int               `json:"raster_count"`
	Type           string            `json:"array_type"`
	XSize          int               `json:"x_size"`
	YSize          int               `json:"y_size"`
	ProjWKT        string            `json:"proj_wkt"`
	Proj4          string            `json:"proj4"`
	GeoTransform   []float64         `json:"geotransform"`
}

type GeoFile struct {
	FileName *string       `json:"filename,omitempty"`
	Driver   string        `json:"file_type"`
	DataSets []GeoMetaData `json:"geo_metadata"`
}

var parserStrings map[string]string = map[string]string{"landsat": `LC(?P<mission>\d)(?P<path>\d\d\d)(?P<row>\d\d\d)(?P<year>\d\d\d\d)(?P<julian_day>\d\d\d)(?P<processing_level>[a-zA-Z0-9]+)_(?P<band>[a-zA-Z0-9]+)`,
	"modis43A4":     `^MCD43A4.A(?P<year>\d\d\d\d)(?P<julian_day>\d\d\d).(?P<horizontal>h\d\d)(?P<vertical>v\d\d).(?P<resolution>\d\d\d).[0-9]+`,
	"modis1":        `^(?P<product>MCD\d\d[A-Z]\d).A(?P<year>\d\d\d\d)(?P<julian_day>\d\d\d).(?P<horizontal>h\d\d)(?P<vertical>v\d\d).(?P<resolution>\d\d\d).[0-9]+`,
	"modis2":        `M(?P<satellite>[OD|YD])(?P<product>[0-9]+_[A-Z0-9]+).A[0-9]+.[0-9]+.(?P<collection_version>\d\d\d).(?P<year>\d\d\d\d)(?P<julian_day>\d\d\d)(?P<hour>\d\d)(?P<minute>\d\d)(?P<second>\d\d)`,
	"modisJP":       `^(?P<product>FC).v302.(?P<root_product>MCD\d\d[A-Z]\d).h(?P<horizontal>\d\d)v(?P<vertical>\d\d).(?P<year>\d\d\d\d).(?P<resolution>\d\d\d).`,
	"modisJP_LR":    `^(?P<product>FC_LR).v302.(?P<root_product>MCD\d\d[A-Z]\d).h(?P<horizontal>\d\d)v(?P<vertical>\d\d).(?P<year>\d\d\d\d).(?P<resolution>\d\d\d).`,
	"himawari8":     `^(?P<year>\d\d\d\d)(?P<month>\d\d)(?P<day>\d\d)(?P<hour>\d\d)(?P<minute>\d\d)(?P<second>\d\d)-P1S-(?P<product>ABOM[0-9A-Z_]+)-PRJ_GEOS141_(?P<resolution>\d+)-HIMAWARI8-AHI`,
	"agdc_landsat1": `LS(?P<mission>\d)_(?P<sensor>[A-Z]+)_(?P<correction>[A-Z]+)_(?P<epsg>\d+)_(?P<x_coord>-?\d+)_(?P<y_coord>-?\d+)_(?P<year>\d\d\d\d).`,
	"elevation_ga":  `^Elevation_1secSRTM_DEMs_v1.0_DEM-S_Tiles_e(?P<longitude>\d+)s(?P<latitude>\d+)dems.nc$`,
	"chirps2.0":     `^chirps-v2.0.(?P<year>\d\d\d\d).dekads.nc$`,
	"era-interim":   `^(?P<product>[a-z0-9]+)_3hrs_ERAI_historical_fc-sfc_(?P<start_year>\d\d\d\d)(?P<start_month>\d\d)(?P<start_day>\d\d)_(?P<end_year>\d\d\d\d)(?P<end_month>\d\d)(?P<end_day>\d\d).nc$`,
	"agdc_landsat2": `LS(?P<mission>\d)_OLI_(?P<sensor>[A-Z]+)_(?P<product>[A-Z]+)_(?P<epsg>\d+)_(?P<x_coord>-?\d+)_(?P<y_coord>-?\d+)_(?P<year>\d\d\d\d).`,
	"agdc_dem":      `SRTM_(?P<product>[A-Z]+)_(?P<x_coord>-?\d+)_(?P<y_coord>-?\d+)_(?P<year>\d\d\d\d)(?P<month>\d\d)(?P<day>\d\d)(?P<hour>\d\d)(?P<minute>\d\d)(?P<second>\d\d)`}

var parsers map[string]*regexp.Regexp = map[string]*regexp.Regexp{}

func init() {
	for key, value := range parserStrings {
		parsers[key] = regexp.MustCompile(value)
	}
}

func parseName(filePath string) (map[string]string, time.Time) {

	for _, r := range parsers {
		_, fileName := filepath.Split(filePath)

		if r.MatchString(fileName) {
			match := r.FindStringSubmatch(fileName)

			result := make(map[string]string)
			for i, name := range r.SubexpNames() {
				if i != 0 {
					result[name] = match[i]
				}
			}

			return result, parseTime(result)
		}
	}

	return nil, time.Time{}
}

func parseTime(nameFields map[string]string) time.Time {
	if _, ok := nameFields["year"]; ok {
		year, _ := strconv.Atoi(nameFields["year"])
		t := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
		if _, ok := nameFields["julian_day"]; ok {
			julianDay, _ := strconv.Atoi(nameFields["julian_day"])
			t = t.Add(time.Hour * 24 * time.Duration(julianDay-1))
		}
		if _, ok := nameFields["month"]; ok {
			if _, ok := nameFields["day"]; ok {
				month, _ := strconv.Atoi(nameFields["month"])
				day, _ := strconv.Atoi(nameFields["day"])
				t = time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
			}
		}
		if _, ok := nameFields["hour"]; ok {
			hour, _ := strconv.Atoi(nameFields["hour"])
			t = t.Add(time.Hour * time.Duration(hour))
		}
		if _, ok := nameFields["minute"]; ok {
			minute, _ := strconv.Atoi(nameFields["minute"])
			t = t.Add(time.Minute * time.Duration(minute))
		}
		if _, ok := nameFields["second"]; ok {
			second, _ := strconv.Atoi(nameFields["second"])
			t = t.Add(time.Second * time.Duration(second))
		}

		return t
	}
	return time.Time{}
}

type POSIXDescriptor struct {
	GID   uint32 `json:"gid"`
	Group string `json:"group"`
	UID   uint32 `json:"uid"`
	User  string `json:"user"`
	Size  int64  `json:"size"`
	Mode  string `json:"mode"`
	Type  string `json:"type"`
	INode uint64 `json:"inode"`
	MTime int64  `json:"mtime"`
	ATime int64  `json:"atime"`
	CTime int64  `json:"ctime"`
}

func ExtractPOSIX(path string) string {
	finfo, err := os.Stat(path)
	if err != nil {
		// no such file or dir
		return ""
	}
	var rec syscall.Stat_t
	err = syscall.Lstat(path, &rec)
	if err != nil {
		// no such file or dir
		return ""
	}
	group, err := user.LookupGroupId(fmt.Sprintf("%d", rec.Gid))
	if err != nil {
		// no such file or dir
		return ""
	}
	user, err := user.LookupId(fmt.Sprintf("%d", rec.Uid))
	if err != nil {
		// no such file or dir
		return ""
	}

	descr := POSIXDescriptor{Group: group.Name, User: user.Username, INode: rec.Ino, MTime: rec.Mtim.Sec, CTime: rec.Ctim.Sec, ATime: rec.Atim.Sec, Size: finfo.Size(), Mode: finfo.Mode().String(), Type: "file", UID: rec.Uid, GID: rec.Gid}

	out, _ := json.Marshal(&descr)
	return fmt.Sprintf("%s\tposix\t%s", path, string(out))
}

func main() {

	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		parts := strings.Split(s.Text(), "\t")
		if len(parts) != 2 {
			fmt.Printf("Input not recognised: %s\n", s.Text())
		}

		fmt.Println(ExtractPOSIX(parts[0]))

		gdalFile := geolib.GDALFile{}
		err := json.Unmarshal([]byte(parts[1]), &gdalFile)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		geoFile := GeoFile{Driver: gdalFile.Driver}
		nameFields, timeStamp := parseName(parts[0])

		if nameFields != nil {

			for _, ds := range gdalFile.DataSets {
				if ds.ProjWKT != "" {
					poly := geolib.GetPolygonFromGeoTransform(ds.ProjWKT, ds.GeoTransform, ds.XSize, ds.YSize)
					//polyWGS84 := poly.ReprojectToWGS84()
					//polyClipped := geolib.ClipDateLine(poly)
					//polyClippedWGS84 := polyClipped.ReprojectToWGS84()

					var times []time.Time
					if nc_times, ok := ds.Extras["nc_times"]; ok {
						for _, timestr := range nc_times {
							t, err := time.ParseInLocation("2006-01-02T15:04:05Z", timestr, time.UTC)
							if err != nil {
								fmt.Println(err)
							}
							times = append(times, t)
						}
					} else {
						times = []time.Time{timeStamp}
					}

					geoFile.DataSets = append(geoFile.DataSets, GeoMetaData{DataSetName: ds.DataSetName, TimeStamps: times, FileNameFields: nameFields, Polygon: poly.ToWKT(), RasterCount: ds.RasterCount, Type: ds.Type, XSize: ds.XSize, YSize: ds.YSize, ProjWKT: ds.ProjWKT, Proj4: poly.Proj4(), GeoTransform: ds.GeoTransform})
					//geoFile.DataSets = append(geoFile.DataSets, GeoMetaData{DataSetName: ds.DataSetName, TimeStamps: times, FileNameFields: nameFields, Polygon: poly.ToWKT(), PolygonWGS84: polyWGS84.ToWKT(), PolygonClip: polyClipped.ToWKT(), PolygonClipWGS84: polyClippedWGS84.ToWKT(),RasterCount: ds.RasterCount, Type: ds.Type, XSize: ds.XSize, YSize: ds.YSize, ProjWKT: ds.ProjWKT, Proj4: poly.Proj4(), GeoTransform: ds.GeoTransform})
				}
			}

			out, err := json.Marshal(&geoFile)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			fmt.Printf("%s\tgdal\t%s\n", parts[0], string(out))
		} else {
			log.Printf("%s non parseable", parts[0])
		}
	}
}
