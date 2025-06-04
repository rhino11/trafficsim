module.exports = {
    testEnvironment: 'jsdom',
    testMatch: ['**/web/static/js/**/*.test.js'],
    setupFilesAfterEnv: ['<rootDir>/web/static/js/test-setup.js'],
    collectCoverageFrom: [
        'web/static/js/**/*.js',
        '!web/static/js/**/*.test.js',
        '!web/static/js/test-setup.js'
    ],
    coverageDirectory: 'coverage-js',
    coverageReporters: ['text', 'lcov', 'html']
};
