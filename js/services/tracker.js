services.factory('Tracker', function() {

    var expenses = [],
        incomes = [];

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
         * Add a new income.
         * Takes a numerical value, or an object containing details of the income:
         *   amount: number
         *   comment: string
         *   date: date
         */
        addIncome: function(income) {
            if (!income) { return; }

            if (typeof income === 'object') {
                income.amount = parseFloat(income.amount);
                incomes.push(income);
            } else if (typeof income === 'number') {
                var incomeObj = { amount: parseFloat(income) };
                incomes.push(incomeObj);
            }
        },

        getExpenses: function() {
            return expenses;
        },

        getIncomes: function() {
            return incomes;
        },

        totalExpenses: function() {
            return expenses.reduce(function(prev, current) {
                return prev + current.amount;
            }, 0);
        },

        totalIncomes: function() {
            return incomes.reduce(function(prev, current) {
                return prev + current.amount;
            }, 0);
        },
    };

    return tracker;

});
