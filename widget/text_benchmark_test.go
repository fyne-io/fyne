package widget

import (
	"testing"

	"fyne.io/fyne/v2"
)

const loremIpsum = `Lorem ipsum dolor sit amet, consectetur adipiscing elit. Pellentesque quis consectetur nisi. Suspendisse id interdum felis. Sed egestas eget tellus eu pharetra. Praesent pulvinar sed massa id placerat. Etiam sem libero, semper vitae consequat ut, volutpat id mi. Mauris volutpat pellentesque convallis. Curabitur rutrum venenatis orci nec ornare. Maecenas quis pellentesque neque. Aliquam consectetur dapibus nulla, id maximus odio ultrices ac. Sed luctus at felis sed faucibus. Cras leo augue, congue in velit ut, mattis rhoncus lectus.

Praesent viverra, mauris ut ullamcorper semper, leo urna auctor lectus, vitae vehicula mi leo quis lorem. Nullam condimentum, massa at tempor feugiat, metus enim lobortis velit, eget suscipit eros ipsum quis tellus. Aenean fermentum diam vel felis dictum semper. Duis nisl orci, tincidunt ut leo quis, luctus vehicula diam. Sed velit justo, congue id augue eu, euismod dapibus lacus. Proin sit amet imperdiet sapien. Mauris erat urna, fermentum et quam rhoncus, fringilla consequat ante. Vivamus consectetur molestie odio, ac rutrum erat finibus a. Suspendisse id maximus felis. Sed mauris odio, mattis eget mi eu, consequat tempus purus.

Nulla facilisi. In a condimentum dolor. Class aptent taciti sociosqu ad litora torquent per conubia nostra, per inceptos himenaeos. Maecenas ultrices, justo vel commodo cursus, sapien neque rhoncus turpis, vitae gravida lorem nisl in tortor. Vestibulum et consectetur sem. Morbi dapibus maximus vulputate. Maecenas nunc lacus, posuere ut gravida id, rhoncus ut velit.

Mauris ullamcorper nec leo quis elementum. Nam vel purus consequat, sodales lacus auctor, interdum massa. Vestibulum non efficitur elit. Mauris laoreet nunc ultricies purus condimentum, egestas hendrerit mauris mollis. Aliquam tincidunt eros sed leo eleifend, varius consequat leo volutpat. Integer eget ultricies lorem, ac vulputate velit. Maecenas eu magna mauris.

Nullam eu mattis dolor. Sed sit amet ipsum gravida, pretium justo eget, mattis est. Cras viverra aliquet velit, a faucibus urna luctus vel. Donec vehicula turpis ligula, non auctor justo tempus nec. In libero orci, tempus vitae ante eu, convallis dapibus nulla. Donec egestas volutpat elit vel semper. Lorem ipsum dolor sit amet, consectetur adipiscing elit. Curabitur enim orci, scelerisque ut viverra at, condimentum sed augue. Duis ut nulla dapibus, mattis ex ac, maximus neque. Integer a augue justo. Ut facilisis enim diam, condimentum interdum elit pulvinar tincidunt. Donec et ipsum ac turpis tristique efficitur. Morbi sodales, odio et vehicula eleifend, quam metus commodo lacus, nec gravida lectus augue eget neque. Quisque tortor lectus, viverra non elit et, hendrerit ornare purus. Cras pretium, felis vel gravida fermentum, libero metus dapibus ex, nec mollis ante enim eget dolor.

Nam nibh nulla, ullamcorper id tortor at, lacinia molestie lacus. Sed eget orci ac sem efficitur placerat mattis id arcu. Praesent leo neque, accumsan eget venenatis vel, accumsan at justo. Maecenas venenatis blandit varius. Fusce ut luctus velit. Nullam quis risus enim. Mauris auctor tempor fermentum. Nam nec leo rutrum nibh aliquam ultrices. Nunc ut molestie lacus, vel feugiat tellus. Suspendisse fringilla eu diam eget interdum.

Nunc eget scelerisque nunc. Quisque faucibus mi in magna ullamcorper porttitor. Ut sodales semper est, ut laoreet felis viverra dapibus. Fusce convallis finibus dapibus. Ut tincidunt at urna non imperdiet. Vivamus efficitur lorem a eros cursus, et tempor lectus malesuada. Sed ac lectus quis enim condimentum commodo eu eget lectus.

Pellentesque sollicitudin nisi at sapien egestas aliquet. Nulla sit amet maximus urna, in laoreet odio. In quis leo fringilla, mollis nisl et, aliquam mi. Nulla gravida tempus justo in pulvinar. Vestibulum nunc libero, accumsan id molestie sed, elementum a enim. Vivamus volutpat fermentum risus, quis faucibus nulla laoreet sed. Quisque id posuere velit. Nam blandit non orci eu mattis. Donec dui ipsum, scelerisque ut vulputate sed, cursus vel velit. Maecenas bibendum neque elit, ac placerat ex pulvinar in. Morbi ultricies est ac malesuada hendrerit. Pellentesque et mauris aliquet neque congue ultricies a vitae libero.

Fusce vitae malesuada ipsum, eget viverra nisi. Nullam varius rhoncus pellentesque. In in ligula est. Suspendisse id felis gravida, cursus justo ac, ultrices tellus. Suspendisse potenti. Quisque sit amet enim vitae dolor pulvinar tempus sit amet eget magna. Aliquam condimentum sapien eu lectus feugiat fringilla. Curabitur fringilla, metus id egestas condimentum, massa nibh lacinia velit, eu pulvinar lacus mi sed sem. Nulla volutpat tincidunt lacinia. Nunc tristique ipsum et nulla finibus egestas. Etiam luctus tempor metus a consectetur.

Ut ac pulvinar purus. Pellentesque tellus quam, condimentum at odio id, viverra porttitor tortor. Aliquam vulputate accumsan mattis. Maecenas quis enim in lorem elementum vehicula. Proin accumsan nec enim in pharetra. Phasellus id enim ligula. Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia Curae; Mauris urna sem, euismod id tempor at, consectetur eu sem. Quisque vestibulum hendrerit diam. Duis gravida tristique velit et sodales. Etiam feugiat blandit diam tempor tincidunt.
`

func BenchmarkText_splitLines(b *testing.B) {
	for n := 0; n < b.N; n++ {
		splitLines([]rune(loremIpsum))
	}
}

func benchmarkTextLineBounds(wrap fyne.TextWrap, b *testing.B) {
	text := []rune(loremIpsum)
	textSize := float32(10)
	textStyle := fyne.TextStyle{}
	measurer := func(text []rune) float32 {
		return fyne.MeasureText(string(text), textSize, textStyle).Width
	}
	for n := 0; n < b.N; n++ {
		lineBounds(text, wrap, 10, measurer)
	}
}

func BenchmarkText_lineBounds_WrapOff(b *testing.B) {
	benchmarkTextLineBounds(fyne.TextWrapOff, b)
}

func BenchmarkText_lineBounds_Truncate(b *testing.B) {
	benchmarkTextLineBounds(fyne.TextTruncate, b)
}

func BenchmarkText_lineBounds_WrapBreak(b *testing.B) {
	benchmarkTextLineBounds(fyne.TextWrapBreak, b)
}

func BenchmarkText_lineBounds_WrapWord(b *testing.B) {
	benchmarkTextLineBounds(fyne.TextWrapWord, b)
}
