var setDisplayTime = function(value){
    if(+value < 10){
        value = '0' + value;
    }
    return value;
};

var textTime = 'Сейчас '+ time.day + ' день ' + time.ten_day +
        ' декады ' + time.month + ' месяца ' + time.year +
        ' года, ' + setDisplayTime(+time.hour) + ':' +
        setDisplayTime(+time.minute);
$('.worldDate').text(textTime);

function getPersonState(state) {
	var curState = "", curStateValue = "";
	for(var key in states) {
		if(state == states[key]) {
			curState = key;
		}
	}
	switch(curState) {
		case "sleep":
			curStateValue = "Спит";
			break;
		case "chores":
			curStateValue = "Занимается домашними делами";
			break;
		default:
			curStateValue = "Нужно добавить описание!";
	}
	return curStateValue;
}

function createPerson(person) {
	return $('<div>', {
		class: "person",
		attr: {
			"data-person-id": person.PersonId
		},
		on: {
			click: function(event) {
				if($(this).find('.person-charasteristics').css('display') == 'none') {
	        		$(this).find('.person-charasteristics').animate({height: 'show'}, 200);
	    		} else {
	        		$(this).find('.person-charasteristics').animate({height: 'hide'}, 100);
	    		}
			}
		},
		append: $('<span>', {
			class: "name",
			text: person.Name
		})
		.add($('<div>', {
			class: "person-charasteristics",
			append: $('<div>', {
				class: "state",
				text: "Состояние: " + getPersonState(person.PersonChr.State)
			})
			.add($('<div>', {
				class: "hunger",
				text: "Голод: " + person.PersonChr.Hunger.toFixed(2)
			}))
			.add($('<div>', {
				class: "thirst",
				text: "Жажда: " + person.PersonChr.Thirst.toFixed(2)
			}))
			.add($('<div>', {
				class: "fatigue",
				text: "Усталость: " + person.PersonChr.Fatigue.toFixed(2)
			}))
			.add($('<div>', {
				class: "somnolency",
				text: "Сонливость: " + person.PersonChr.Somnolency.toFixed(2)
			}))
		}))
	});
}

function drawPersons(persons) {
	persons.forEach(function(person, i){
		$('.persons').append(createPerson(person));
	})
}

function updatePersonChr(persons) {
	persons.forEach(function(person, i){
		$('div[data-person-id='+person.PersonId+']').find('.state').text("Состояние: " + getPersonState(person.PersonChr.State));
		$('div[data-person-id='+person.PersonId+']').find('.hunger').text("Голод: " + person.PersonChr.Hunger);
		$('div[data-person-id='+person.PersonId+']').find('.thirst').text("Жажда: " + person.PersonChr.Thirst);
		$('div[data-person-id='+person.PersonId+']').find('.fatigue').text("Усталость: " + person.PersonChr.Fatigue);
		$('div[data-person-id='+person.PersonId+']').find('.somnolency').text("Сонливость: " + person.PersonChr.Somnolency);
	})
}

drawPersons(persons);