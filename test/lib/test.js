/**
 * Created by Than on 09.07.2016.
 */
// Data is read from select statements published by server (further down)
players = new PgSubscription('allPlayers');

// Extra (not used anywhere on the app UI) subscription to display different
//  use case with arguments and manually authored triggers
myScore = new PgSubscription('playerScore', 'Maxwell');

timer = new PgSubscription('worldTime');
