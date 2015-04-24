controllers.controller('TrackingController', ['$location', 'DB', 'Tracker',
  function TrackingControllerFactory($location, DB, Tracker) {

    var track = this;
    var INCOME = 'income',
        EXPENSE = 'expense';

    DB.getExpenses('test@example.com')
      .success(function(data) {
        console.log('expenses: %O', data)
        });

    track.getExpenses = Tracker.getExpenses;
    track.getTotalExpenses = Tracker.totalExpenses;

    track.addExpense = function() {
        var expense = {
            amount: track.amount,
            comment: track.comment,
            date: track.date,
            kind: EXPENSE
        };

        Tracker.addExpense(expense);

        DB.storeExpense(expense)
          .success(function(data, status) {
            if (data.ok) {
                console.log('Stored %s successully!', expense.kind);
                Tracker.saveSuccess(Tracker.findExpense(expense));
            } else {
                console.log('Failed to store expense!');
                Tracker.saveError(Tracker.findExpense(expense));
            }
            console.log(data);
          })
          .error(function(data, status) {
            console.log('Failed to store expense - HTTP %s', status);
            console.log(data);
          });

        track.clearForm();
    };

    track.getIncomes = Tracker.getIncomes;
    track.getTotalIncomes = Tracker.totalIncomes;

    track.addIncome = function() {
        var income = {
            amount: track.amount,
            comment: track.comment,
            date: track.date,
            kind: INCOME
        };

        Tracker.addIncome(income);

        DB.storeIncome(income)
          .success(function(data, status) {
            if (data.ok) {
                console.log('Stored %s successully!', income.kind);
                Tracker.saveSuccess(Tracker.findIncome(income));
            } else {
                console.log('Failed to store income!');
                Tracker.saveError(Tracker.findIncome(income));
            }
            console.log(data);
          })
          .error(function(data, status) {
            console.log('Failed to store income - HTTP %s', status);
            console.log(data);
          });

        track.clearForm();
    };

    track.clearForm = function() {
        track.amount = null;
        track.comment = null;
        track.date = new Date();
    };

    track.clearForm();
}]);
