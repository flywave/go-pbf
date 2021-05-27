package pbf

import "math"

var powerfactor = math.Pow(10.0, 7.0)

func (pbf *Reader) ReadSVarintPower() float64 {
	num := int(pbf.ReadVarint())
	if num%2 == 1 {
		return float64((num+1)/-2) / powerfactor
	} else {
		return float64(num/2) / powerfactor
	}
}

func (pbf *Reader) ReadPoint(endpos int) []float64 {
	for pbf.Pos < endpos {
		x := pbf.ReadSVarintPower()
		y := pbf.ReadSVarintPower()
		return []float64{Round(x, .5, 7), Round(y, .5, 7)}
	}
	return []float64{}
}

func (pbf *Reader) ReadLine(num int, endpos int) [][]float64 {
	var x, y float64
	if num == 0 {

		for startpos := pbf.Pos; startpos < endpos; startpos++ {
			if pbf.Pbf[startpos] <= 127 {
				num += 1
			}
		}
		newlist := make([][]float64, num/2)

		for i := 0; i < num/2; i++ {
			x += pbf.ReadSVarintPower()
			y += pbf.ReadSVarintPower()
			newlist[i] = []float64{Round(x, .5, 7), Round(y, .5, 7)}
		}

		return newlist
	} else {
		newlist := make([][]float64, num/2)

		for i := 0; i < num/2; i++ {
			x += pbf.ReadSVarintPower()
			y += pbf.ReadSVarintPower()

			newlist[i] = []float64{Round(x, .5, 7), Round(y, .5, 7)}

		}
		return newlist
	}
}

func (pbf *Reader) ReadPolygon(endpos int) [][][]float64 {
	polygon := [][][]float64{}
	for pbf.Pos < endpos {
		num := pbf.ReadVarint()
		polygon = append(polygon, pbf.ReadLine(num, endpos))
	}
	return polygon
}

func (pbf *Reader) ReadMultiPolygon(endpos int) [][][][]float64 {
	multipolygon := [][][][]float64{}
	for pbf.Pos < endpos {
		num_rings := pbf.ReadVarint()
		polygon := make([][][]float64, num_rings)
		for i := 0; i < num_rings; i++ {
			num := pbf.ReadVarint()
			polygon[i] = pbf.ReadLine(num, endpos)
		}
		multipolygon = append(multipolygon, polygon)
	}
	return multipolygon
}

func (pbf *Reader) ReadBoundingBox() []float64 {
	bb := make([]float64, 4)
	pbf.ReadVarint()
	bb[0] = float64(pbf.ReadSVarintPower())
	bb[1] = float64(pbf.ReadSVarintPower())
	bb[2] = float64(pbf.ReadSVarintPower())
	bb[3] = float64(pbf.ReadSVarintPower())
	return bb
}
