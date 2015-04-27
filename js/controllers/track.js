controllers.controller('TrackingController', ['$location', 'DB', 'Tracker', 'Type',
  function TrackingControllerFactory($location, DB, Tracker, Type) {

    var track = this;

    track.getExpenses = Tracker.getExpenses;
    track.getTotalExpenses = Tracker.totalExpenses;
    track.getIncomes = Tracker.getIncomes;
    track.getTotalIncomes = Tracker.totalIncomes;

    track.setDefaultCategories = setDefaultCategories.bind(track);
    track.categories = [];
    track.setDefaultCategories();

    track.addExpense = function() {
        var expense = {
            amount: track.amount,
            category: track.category,
            comment: track.comment,
            date: track.date,
            kind: Type.EXPENSE
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

    track.addIncome = function() {
        var income = {
            amount: track.amount,
            category: track.category,
            comment: track.comment,
            date: track.date,
            kind: Type.INCOME
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
        track.category = null;
        track.comment = null;
        track.date = new Date();
    };

    track.clearForm();

    function setDefaultCategories() {
      this.categories = [
        'Auto',
        'Credit Card',
        'Clothing',
        'Utilities',
        'Entertainment',
        'Transportation',
        'Grocery',
        'Health',
        'Home',
        'Other'
      ];
    }
}]);
