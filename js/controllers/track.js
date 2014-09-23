controllers.controller('TrackingController', ['$location', 'Tracker', function TrackingControllerFactory($location, Tracker) {

    var track = this;

    track.getExpenses = Tracker.getExpenses;
    track.getTotalExpenses = Tracker.totalExpenses;

    track.addExpense = function() {
        var expense = {
            amount: track.newExpenseAmount,
            comment: track.newExpenseComment,
            date: track.newExpenseDate,
        };
        Tracker.addExpense(expense);
        track.clearExpenseEntry();
    };

    track.getCredits = Tracker.getCredits;
    track.getTotalCredits = Tracker.totalCredits;

    track.addCredit = function() {
        var credit = {
            amount: track.newCreditAmount,
            comment: track.newCreditComment,
            date: track.newCreditDate,
        };
        Tracker.addCredit(credit);
        track.clearCreditEntry();
    };
    
    track.clearExpenseEntry = function() {
        track.newExpenseAmount = null;
        track.newExpenseComment = null;
        track.newExpenseDate = null;
    };

    track.clearCreditEntry = function() {
        track.newCreditAmount = null;
        track.newCreditComment = null;
        track.newCreditDate = null;
    };

}]);
