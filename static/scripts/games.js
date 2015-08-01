"use strict";

var colors = [
	"#16a085", 
	"#8e44ad", 
	"#c0392b"
];
var numbers = [6,8,9];

function game0() {
	var lines = game0generator(8, 25);
	var list = $('<ul>').addClass('game0');

	$.each(lines, function(i) {
	    var li = $('<li/>')
	    	.attr('style', 'height: 35px')
	        .addClass('no-list')
	        .appendTo(list);
	    var subList = $('<ul>').appendTo(li);

	    numbers = lines[i]
	    $.each(numbers, function(j) {
	    	var num = numbers[j];
	    	var elem = $('<li/>')
	    		.attr('style', 'cursor: pointer; display: inline; color: ' + num.color)
	    		.addClass('no-list')
	    		.text(num.val)
	    		.appendTo(subList);
	    });
	});

console.log(lines);
console.log(list);
console.log(list[0].outerHTML);

	return list[0].outerHTML;
}

function game0generator(linesNum, lineLen) {
	var mainNum = getRandomFromArray(numbers);
	var mainColor = getRandomFromArray(colors);

	var lines = [];
	for (var i = 0; i < linesNum; i++) {
		var line = [];
		for (var j = 0; j < lineLen; j++) {
			var val = mainNum
			if (Math.random() < 2 / lineLen) {
				val = getRandomFromArray(numbers);
			}
			var number = {
				val: val,
				color: mainColor
			};
			line.push(number);
		}
		lines.push(line)
	}
	
	return lines;
};


function getRandomFromArray(arr) {
	return arr[Math.floor(Math.random() * arr.length)];
}