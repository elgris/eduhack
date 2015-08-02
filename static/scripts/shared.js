var inProgress = false;

function score (success, value) {
	value = value || 1;
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

function getRandomFromArray(arr) {
	return arr[Math.floor(Math.random() * arr.length)];
}

function finish (message, success) {
	if (!inProgress) {
		return
	}
    $("#dialog-confirm").html(message);
    var title = "";
    if (success) {
        title = "Congratulations!";
        score(true, 10)
    } else {
        title = "Let's solve next problem!";
    }
	$("#dialog-confirm").dialog({
        resizable: false,
        modal: true,
        title: title,
        height: 250,
        width: 400,
        buttons: {
            "OK": function () {
            	socket.emit('finish', message);
                $(this).closest('.ui-dialog-content').dialog('close'); 
            }
        }
    });


}

function notify (str) {
	var ev = {
		'score_player_1': $player1score.html(),
		'score_player_2': $player2score.html(),
		'message': str
	};

	socket.emit('event', JSON.stringify(ev));
}

function hideTile (line, col) {
	var $elem = $('ul.game li:nth-child(' + (line+1) + ') ul li:nth-child(' + (col+1) + ')')
	var style = 'visibility: hidden;' + $elem.attr('style');
	$elem.attr('style', style);
}