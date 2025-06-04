// Mock Leaflet globally
global.L = {
    map: jest.fn(() => ({
        addLayer: jest.fn(),
        removeLayer: jest.fn(),
        setView: jest.fn(),
        getZoom: jest.fn(() => 10),
        getBounds: jest.fn(() => ({
            contains: jest.fn(() => true),
            getNorth: jest.fn(() => 90),
            getSouth: jest.fn(() => -90),
            getEast: jest.fn(() => 180),
            getWest: jest.fn(() => -180)
        })),
        on: jest.fn(),
        off: jest.fn(),
        panTo: jest.fn(),
        fitBounds: jest.fn()
    })),
    layerGroup: jest.fn(() => ({
        addTo: jest.fn(),
        addLayer: jest.fn(),
        removeLayer: jest.fn(),
        clearLayers: jest.fn(),
        eachLayer: jest.fn(),
        getLayers: jest.fn(() => [])
    })),
    marker: jest.fn(() => ({
        addTo: jest.fn(),
        setLatLng: jest.fn(),
        bindPopup: jest.fn(),
        openPopup: jest.fn(),
        closePopup: jest.fn(),
        remove: jest.fn(),
        on: jest.fn(),
        off: jest.fn()
    })),
    divIcon: jest.fn(() => ({})),
    icon: jest.fn(() => ({})),
    popup: jest.fn(() => ({
        setContent: jest.fn(),
        openOn: jest.fn()
    })),
    polyline: jest.fn(() => ({
        addTo: jest.fn(),
        remove: jest.fn(),
        setLatLngs: jest.fn()
    })),
    latLng: jest.fn((lat, lng) => ({ lat, lng })),
    latLngBounds: jest.fn(() => ({
        contains: jest.fn(() => true),
        extend: jest.fn()
    })),
    canvas: jest.fn(() => ({
        drawing: true,
        bringToFront: jest.fn(),
        bringToBack: jest.fn()
    })),
    svg: jest.fn(() => ({
        drawing: false,
        bringToFront: jest.fn(),
        bringToBack: jest.fn()
    })),
    tileLayer: jest.fn(() => ({
        addTo: jest.fn(),
        remove: jest.fn()
    })),
    control: {
        scale: jest.fn(() => ({
            addTo: jest.fn()
        })),
        zoom: jest.fn(() => ({
            addTo: jest.fn()
        })),
        attribution: jest.fn(() => ({
            addTo: jest.fn()
        }))
    }
};

// Mock Canvas API
global.HTMLCanvasElement.prototype.getContext = jest.fn(() => ({
    fillRect: jest.fn(),
    clearRect: jest.fn(),
    getImageData: jest.fn(),
    putImageData: jest.fn(),
    createImageData: jest.fn(),
    setTransform: jest.fn(),
    drawImage: jest.fn(),
    save: jest.fn(),
    fillText: jest.fn(),
    restore: jest.fn(),
    beginPath: jest.fn(),
    moveTo: jest.fn(),
    lineTo: jest.fn(),
    closePath: jest.fn(),
    stroke: jest.fn(),
    translate: jest.fn(),
    scale: jest.fn(),
    rotate: jest.fn(),
    arc: jest.fn(),
    fill: jest.fn(),
    measureText: jest.fn(() => ({ width: 0 })),
    transform: jest.fn(),
    rect: jest.fn(),
    clip: jest.fn()
}));

// Mock RequestAnimationFrame
global.requestAnimationFrame = jest.fn(cb => setTimeout(cb, 0));
global.cancelAnimationFrame = jest.fn(id => clearTimeout(id));

// Mock Performance API
global.performance = {
    now: jest.fn(() => Date.now()),
    mark: jest.fn(),
    measure: jest.fn(),
    getEntriesByName: jest.fn(() => [])
};

// Mock WebSocket
global.WebSocket = jest.fn(() => ({
    send: jest.fn(),
    close: jest.fn(),
    addEventListener: jest.fn(),
    removeEventListener: jest.fn(),
    readyState: 1, // OPEN
    CONNECTING: 0,
    OPEN: 1,
    CLOSING: 2,
    CLOSED: 3
}));

// Mock console methods
global.console = {
    ...console,
    log: jest.fn(),
    error: jest.fn(),
    warn: jest.fn(),
    info: jest.fn(),
    debug: jest.fn()
};
