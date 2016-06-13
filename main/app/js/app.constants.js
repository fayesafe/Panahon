angular
  .module('weatherApp')
  .constant('EVENTS', {
    'DATA_UPDATED': 'data_updated'
  })
  .constant('CHART_OPTIONS', {
    'TEMPERATURE': {
      data: {
        json: [],
        keys: {
          x: 'time',
          value: ['temperature']
        },
        names: { temperature: 'Temperatur' },
        classes: { temperature: 'temperature' },
        types: { temperature: 'spline' },
        labels: {
          format: {
            temperature: function (v, id, i, j) { return v + '°C'; }
          }
        },
        colors: { temperature: 'orange' }
      },
      axis: {
        x: {
          type: 'timeseries',
          tick: { format: '%H:%M:%S' },
          labels: true
        }
      },
      grid: { y: { lines: [{ value: 0 }] } }
    },
    'HUMIDITY': {
      data: {
        json: [],
        keys: {
          x: 'time',
          value: ['humidity']
        },
        names: { humidity: 'Luftfeuchtigkeit' },
        classes: { humidity: 'humidity' },
        types: { humidity: 'line' },
        labels: {
          format: {
            humidity: function (v, id, i, j) { return v + '%'; }
          }
        },
        colors: { humidity: 'green' }
      },
      axis: {
        x: {
          type: 'timeseries',
          tick: { format: '%H:%M:%S' }
        }
      }
    },
    'PRESSURE': {
      data: {
        json: [],
        keys: {
          x: 'time',
          value: ['pressure']
        },
        names: {pressure: 'Luftdruck' },
        classes: { pressure: 'pressure' },
        type: 'bar',
        labels: {
          format: function (v, id, i, j) { return v + 'hPa'; }
        },
        colors: { pressure: 'blue' },
        empty: {
          label: {
            text: "Loading data..."
          }
        }
      },
      axis: {
        x: {
          type: 'timeseries',
          tick: { format: '%H:%M:%S' },
        },
        y: {
          tick: { format: function (d) {
            return d + 'hPa';
          } },
          label: {
            text: 'Your Y Axis',
            position: 'outer-middle',
          }
        }
      }
    }
  })
  .run(function ($rootScope, CHART_OPTIONS) {
    $rootScope.CHART_OPTIONS = CHART_OPTIONS;
  });
