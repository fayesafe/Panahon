angular
  .module('weatherApp')
  .constant('EVENTS', {
    'DATA_UPDATED': 'data_updated'
  })
  .constant('CHART_OPTIONS', {
    'TEMPERATURE': {
      data: {
        x: 'time',
        rows: [],
        hide: ['humidity', 'sun', 'rain', 'pressure'],
        names: { temperature: 'Temperatur', 'min_temperature': 'Min. Temperatur', 'max_temperature': 'Max. Temperatur' },
        classes: { temperature: 'temperature', 'min_temperature': 'min-temperature', 'max_temperature': 'max-temperature' },
        types: { temperature: 'line', 'min_temperature': 'line', 'max_temperature': 'line' },
        labels: {
          format: {
            temperature: function (v, id, i, j) { return v + '°C'; }
          }
        },
        colors: { temperature: '#ff851b',  'min_temperature': '#3498db', 'max_temperature': '#e74c3c' },
        empty: {
          label: {
            text: "Loading data..."
          }
        }
      },
      legend: {
        hide: ['humidity', 'sun', 'rain', 'pressure']
      },
      axis: {
        x: {
          type: 'timeseries',
          tick: { format: '%H:%M' },
          labels: true
        }
      },
      grid: { y: { lines: [{ value: 0 }] } }
    },
    'HUMIDITY': {
      data: {
        x: 'time',
        rows: [],
        hide: ['temperature', 'sun', 'rain', 'pressure', 'min_temperature', 'max_temperature'],
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
      legend: {
        hide: ['temperature', 'sun', 'rain', 'pressure', 'min_temperature', 'max_temperature']
      },
      axis: {
        x: {
          type: 'timeseries',
          tick: { format: '%H:%M' }
        }
      }
    },
    'PRESSURE': {
      data: {
        x: 'time',
        rows: [],
        hide: ['humidity', 'sun', 'rain', 'temperature', 'min_temperature', 'max_temperature'],
        names: { pressure: 'Luftdruck' },
        classes: { pressure: 'pressure' },
        type: 'bar',
        colors: { pressure: '#158cba' },
        empty: {
          label: {
            text: "Loading data..."
          }
        }
      },
      legend: {
        hide: ['humidity', 'sun', 'rain', 'temperature', 'min_temperature', 'max_temperature']
      },
      axis: {
        x: {
          type: 'timeseries',
          tick: { format: '%H:%M' },
        },
        y: {
          tick: { format: function (d) {
            return d + 'hPa';
          } }
        }
      }
    },
    'RAIN': {
      data: {
        x: 'time',
        rows: [],
        hide: ['humidity', 'sun', 'temperature', 'pressure', 'min_temperature', 'max_temperature'],
        names: { rain: 'Regen' },
        classes: { rain: 'rain' },
        types: { rain: 'area-step' },
        colors: { rain: '#158cba' },
        empty: {
          label: {
            text: "Loading data..."
          }
        }
      },
      legend: {
        hide: ['time', 'humidity', 'sun', 'temperature', 'pressure', 'min_temperature', 'max_temperature']
      },
      axis: {
        x: {
          type: 'timeseries',
          tick: { format: '%H:%M' },
        },
        y: {
          tick: { format: function (d) {
            return d + '%';
          } }
        }
      }
    },
    'SUN': {
      data: {
        x: 'time',
        rows: [],
        hide: ['humidity', 'rain', 'temperature', 'pressure', 'min_temperature', 'max_temperature'],
        names: { sun: 'Sonne' },
        classes: { sun: 'rain' },
        types: { sun: 'spline' },
        colors: { sun: '#ff4136' },
        empty: {
          label: {
            text: "Loading data..."
          }
        }
      },
      legend: {
        hide: ['time', 'humidity', 'rain', 'temperature', 'pressure', 'min_temperature', 'max_temperature']
      },
      axis: {
        x: {
          type: 'timeseries',
          tick: { format: '%H:%M' },
        },
        y: {
          inverted: true,
          label: {
            text: 'Helligkeit',
            position: 'outer-middle'
          },
          tick: { format: function (d) { return ''; } }
        }
      }
    },
    'REALTIME': {
      size: {
        height: 500
      },
      data: {
        x: 'time',
        rows: [],
        axes: {
          ts: 'x',
          temperature: 'y',
          pressure: 'y2',
          sun: 'y2'
        },
        names: { temperature: 'Temperatur', humidity: 'Luftfeuchtigkeit', pressure: 'Luftdruck', rain: 'Regen', sun: 'Sonne', 'min_temperature': 'Min. Temperatur', 'max_temperature': 'Max. Temperatur'  },
        classes: { temperature: 'temperature', 'min_temperature': 'min-temperature', 'max_temperature': 'max-temperature' },
        types: { temperature: 'spline', 'min_temperature': 'spline', 'max_temperature': 'spline', humidity: 'line', pressure: 'bar', rain: 'area-step', sun: 'spline' },
        labels: {
          format: {
            temperature: function (v, id, i, j) { return v.toFixed(1) + '°C'; },
            min_temperature: function (v, id, i, j) { return v.toFixed(1) + '°C'; },
            max_temperature: function (v, id, i, j) { return v.toFixed(1) + '°C'; },
            humidity: function (v, id, i, j) { return v.toFixed(1) + '%'; },
          }
        },
        colors: { temperature: '#ff851b', 'min_temperature': '#3498db', 'max_temperature': '#e74c3c', humidity: '#28b62c', pressure: '#C5EFF7', sun: '#ff4136' },
        empty: { label: { text: "Loading data..." } }
      },
      axis: {
        x: {
          type: 'timeseries',
          tick: { format: '%H:%M.%S' },
        },
        y2: {
          show: true,
          tick: { format: function (d) { return d + ' hPa'; } }
        }
      },
      grid: { y: { lines: [{ value: 0 }] } }
    },
    'ALL': {
      size: {
        height: 500
      },
      data: {
        x: 'time',
        rows: [],
        axes: {
          ts: 'x',
          temperature: 'y',
          pressure: 'y2',
          sun: 'y2'
        },
        names: { temperature: 'Temperatur', humidity: 'Luftfeuchtigkeit', pressure: 'Luftdruck', rain: 'Regen', sun: 'Sonne', 'min_temperature': 'Min. Temperatur', 'max_temperature': 'Max. Temperatur'  },
        classes: { temperature: 'temperature', 'min_temperature': 'min-temperature', 'max_temperature': 'max-temperature' },
        types: { temperature: 'spline', 'min_temperature': 'spline', 'max_temperature': 'spline', humidity: 'line', pressure: 'bar', rain: 'area-step', sun: 'spline' },
        labels: {
          format: {
            temperature: function (v, id, i, j) { return v.toFixed(1) + '°C'; },
            min_temperature: function (v, id, i, j) { return v.toFixed(1) + '°C'; },
            max_temperature: function (v, id, i, j) { return v.toFixed(1) + '°C'; },
            humidity: function (v, id, i, j) { return v.toFixed(1) + '%'; },
          }
        },
        colors: { temperature: '#ff851b', 'min_temperature': '#3498db', 'max_temperature': '#e74c3c', humidity: '#28b62c', pressure: '#C5EFF7', sun: '#ff4136' },
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
  .run(function ($rootScope, CHART_OPTIONS) {
    $rootScope.CHART_OPTIONS = CHART_OPTIONS;
  });
