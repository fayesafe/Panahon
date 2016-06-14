angular
  .module('weatherApp')
  .constant('EVENTS', {
    'DATA_UPDATED': 'data_updated'
  })
  .constant('DATETIME', {
    'MONTHS': [
      'Jan', 'Febr', 'März', 'Apr', 'Mai', 'Juni', 'Juli', 'Aug', 'Sep', 'Okt', 'Nov', 'Dez'
    ]
  })
  .constant('CHART_OPTIONS', {
    'TEMPERATURE': {
      data: {
        json: [],
        keys: {
          x: 'ts',
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
        colors: { temperature: '#ff851b' },
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
          labels: true
        }
      },
      grid: { y: { lines: [{ value: 0 }] } }
    },
    'HUMIDITY': {
      data: {
        json: [],
        keys: {
          x: 'ts',
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
        colors: { humidity: '#28b62c' },
        empty: {
          label: {
            text: "Loading data..."
          }
        }
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
          x: 'ts',
          value: ['pressure']
        },
        names: { pressure: 'Luftdruck' },
        classes: { pressure: 'pressure' },
        type: 'bar',
        labels: {
          format: function (v, id, i, j) { return v + 'hPa'; }
        },
        colors: { pressure: '#158cba' },
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
    },
    'ALL': {
      bindto: '#chart-all',
      data: {
        json: [],
        axis: {
          ts: 'x',
          temperature: 'y',
          pressure: 'y2'
        },
        keys: {
          x: 'ts',
          value: ['temperature', 'humidity', 'pressure']
        },
        axes: {
          'pressure': 'y2'
        },
        names: { temperature: 'Temperatur', humidity: 'Luftfeuchtigkeit', pressure: 'Luftdruck' },
        classes: {  },
        types: { temperature: 'spline', humidity: 'line', pressure: 'bar' },
        labels: {
          format: {
            temperature: function (v, id, i, j) { return v + '°C'; },
            humidity: function (v, id, i, j) { return v + '%'; },
            //pressure: function (v, id, i, j) { return v + 'hPa'; },
          }
        },
        colors: { temperature: '#ff851b', humidity: '#28b62c', pressure: '#3498db' },
        empty: { label: { text: "Loading data..." } }
      },
      axis: {
        x: {
          type: 'category',
          labels: true
        },
        y2: {
          show: true,
          tick: { format: function (d) { return d + ' hPa'; } }
        }
      },
      grid: { y: { lines: [{ value: 0 }] } },
      zoom: {
          enabled: true
      },
      subchart: {
          show: true
      }
    }
  })
  .run(function ($rootScope, CHART_OPTIONS, DATETIME) {
    $rootScope.CHART_OPTIONS = CHART_OPTIONS;
    $rootScope.DATETIME = DATETIME;
  });
