/**
 * Created by Than on 31.07.2016.
 */
Template.person.helpers({
    characteristics: function () {
        var id = this.id;
        Characteristics.depend();
        //console.log(Characteristics);
        var ch = Characteristics.filter(function(item){
            //console.log(item);
            return item._index == id;
        });
        return ch.length && ch;
    }
});

Template.person.events({
    'click': function(event, template) {
        if(template.$('.person-charasteristics').css('display') == 'none') {
            template.$('.person-charasteristics').animate({height: 'show'}, 200);
        } else {
            template.$('.person-charasteristics').animate({height: 'hide'}, 100);
        }
    }
});