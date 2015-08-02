"use strict";

var colors = [
	"#16a085", 
	"#8e44ad", 
	"#c0392b"
];

var game0data = [];
var game0remains = 0;

function game0(first, second, game0data) {
	if (playerIndex == 1) {
		playerValue = first;
	} else {
		playerValue = second;
	}

	var mainColor = getRandomFromArray(colors);

	var lines = game0data;
	var list = $('<ul>').addClass('game0');

	$.each(lines, function(i) {
	    var li = $('<li/>')
	    	.attr('style', 'height: 35px')
	        .addClass('no-list')
	        .appendTo(list);
	    var subList = $('<ul>').appendTo(li);

	    var numbers = lines[i]
	    $.each(numbers, function(j) {
	    	var num = numbers[j];
	    	if (num == playerValue) {
	    		game0remains++;
	    	}
	    	$('<li onclick="numClicked(' + i + ',' + j + ')"/>')	    		
    			.attr('style', 'cursor: pointer; display: inline; color: ' + mainColor)
	    		.addClass('no-list')
	    		.text(num)
	    		.appendTo(subList);
	    });
	});
	return list[0].outerHTML;
}

function numClicked (line, col) {
	if (game0remains <= 0) {
		return;
	}
	hideNum(line, col)
	var $elem = $('ul.game0 li:nth-child(' + (line+1) + ') ul li:nth-child(' + (col+1) + ')');

	var val = parseInt($elem.html());
	if (val == playerValue) {
		game0remains--;
	}
	score(val == playerValue);
	notify('hideNum(' + line + ',' + col + ')');

	if (game0remains <= 0) {
		finish('You have finished your part!', true);
	}
}

function hideNum (line, col) {
	var $elem = $('ul.game0 li:nth-child(' + (line+1) + ') ul li:nth-child(' + (col+1) + ')')
	var style = 'visibility: hidden;' + $elem.attr('style');
	$elem.attr('style', style);
}

function score (success) {
	$elem = $('#player-' + playerIndex + '-score');
	var val = parseInt($elem.html())
	if (success) {
		setHtml($elem, val + 1);
	} else {
		setHtml($elem, val - 1);
	}
}

function setHtml ($elem, value) {
  $elem.html(value)
  var score1 = parseInt($player1score.html());
  var score2 = parseInt($player2score.html());
  $score.html(score1 + score2);
}


function finish (message, success) {
    $("#dialog-confirm").html(message);
    var title = "";
    if (success) {
        title = "Congratulations!";
    } else {
        title = "You failed :(";
    }
	$("#dialog-confirm").dialog({
        resizable: false,
        modal: true,
        title: title,
        height: 250,
        width: 400,
        buttons: {
            "OK": function () {
                $(this).closest('.ui-dialog-content').dialog('close'); 
            }
        }
    });

    socket.emit('finish', message);
}

function notify (str) {
	var ev = {
		'score_player_1': $player1score.html(),
		'score_player_2': $player2score.html(),
		'message': str
	};

	socket.emit('event', JSON.stringify(ev));
}

function game0generator(linesNum, lineLen, main, first, second) {
	var mainColor = getRandomFromArray(colors);

	var lines = [];
	for (var i = 0; i < linesNum; i++) {
		var line = [];
		for (var j = 0; j < lineLen; j++) {
			var val = main
			var rnd = Math.random()
			if (rnd < 2 / lineLen) {
				val = first;
			} else if (rnd < 4 / lineLen) {
				val = second;
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