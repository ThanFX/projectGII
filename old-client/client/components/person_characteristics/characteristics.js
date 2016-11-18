/**
 * Created by Than on 31.07.2016.
 */
Template.characteristic.helpers({
    hunger: function () {
        return (+this.hunger).toFixed(2);
    },
    thirst: function () {
        return (+this.thirst).toFixed(2);
    },
    fatigue: function () {
        return (+this.fatigue).toFixed(2);
    },
    somnolency: function () {
        return (+this.somnolency).toFixed(2);
    },
    state: function() {
        let s;
        switch (this.state) {
            case 1:
                s = "Спит";
                break;
            case 5:
                s = "Занимается домашними делами";
                break;
        }
        return s;
    }
});