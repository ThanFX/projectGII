/**
 * Created by Than on 09.08.2016.
 */
Template.worldMap.helpers({

});

Template.worldMap.onRendered(() => {
    this.$('.chunk').css({
        'backgroundImage': 'url(resources/forest.png)'
    });
});