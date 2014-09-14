controllers.controller('TrackingController', ['$location', function TrackingControllerFactory($location) {

    var track = this;
    
    track.expenses = [];
    track.credits = [];

    track.expenses.push({
        amount: 4.50,
        comments: 'Bought some ice cream.',
        date: new Date()
    });

}]);
