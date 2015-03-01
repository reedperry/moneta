services.factory('Tracker', function() {

    var expenses = [],
        credits = [];

    var tracker = {

        /*
         * Add a new expense.
         * Takes a numerical value, or an object containing details of the expense:
         *   amount: number
         *   comment: string
         *   date: date
         */
        addExpense: function(expense) {
            if (!expense) { return; }

            if (typeof expense === 'object') {
                expense.amount = parseFloat(expense.amount);
                expenses.push(expense);
            } else if (typeof expense === 'number') {
                var expenseObj = { amount: parseFloat(expense) };
                expenses.push(expenseObj);
            }
        },

        /*
         * Add a new credit.
         * Takes a numerical value, or an object containing details of the credit:
         *   amount: number
         *   comment: string
         *   date: date
         */
        addCredit: function(credit) {
            if (!credit) { return; }

            if (typeof credit === 'object') {
                credit.amount = parseFloat(credit.amount);
                credits.push(credit);
            } else if (typeof credit === 'number') {
                var creditObj = { amount: parseFloat(credit) };
                credits.push(creditObj);
            }
        },

        getExpenses: function() {
            return expenses;
        },

        getCredits: function() {
            return credits;
        },

        totalExpenses: function() {
            return expenses.reduce(function(prev, current) {
                return prev + current.amount;
            }, 0);
        },

        totalCredits: function() {
            return credits.reduce(function(prev, current) {
                return prev + current.amount;
            }, 0);
        },
    };

    return tracker;

});
