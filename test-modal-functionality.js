// Test script for enhanced platform editing modal functionality
// Run this in the browser console when on the scenario-builder page

console.log('=== Enhanced Platform Editing Modal Test ===');

// Test 1: Check if the ScenarioBuilder class is loaded
if (typeof scenarioBuilder !== 'undefined') {
    console.log('✅ ScenarioBuilder instance found');
} else {
    console.log('❌ ScenarioBuilder instance not found');
    return;
}

// Test 2: Check if platforms are loaded
if (scenarioBuilder.platforms && scenarioBuilder.platforms.length > 0) {
    console.log(`✅ Platforms loaded: ${scenarioBuilder.platforms.length} platforms`);
} else {
    console.log('❌ No platforms loaded');
    return;
}

// Test 3: Test enhanced platform data generation
const testPlatform = scenarioBuilder.platforms[0];
console.log('Testing enhanced platform data generation with:', testPlatform.name);

try {
    const enhancedData = scenarioBuilder.enhancePlatformData(testPlatform);
    console.log('✅ Enhanced platform data generated successfully');
    console.log('Enhanced data structure:', {
        name: enhancedData.name,
        milStdData: enhancedData.milStdData,
        specifications: enhancedData.specifications,
        operational: enhancedData.operational
    });
} catch (error) {
    console.log('❌ Error generating enhanced platform data:', error);
}

// Test 4: Test enhanced platform display generation
try {
    const displayHTML = scenarioBuilder.createEnhancedPlatformDisplay(testPlatform);
    console.log('✅ Enhanced platform display HTML generated successfully');
    console.log('Display HTML length:', displayHTML.length, 'characters');
} catch (error) {
    console.log('❌ Error generating enhanced platform display:', error);
}

// Test 5: Check if modal HTML elements exist
const modalElements = [
    'platformEditModal',
    'editPlatformName',
    'editPlatformClass',
    'editSymbolCode',
    'editUnitDesignation',
    'editHigherFormation',
    'editLength',
    'editMaxSpeed',
    'editServiceCeiling',
    'editRange',
    'editCrewSize',
    'editFuelLevel',
    'editMissionReadiness',
    'editMaintenanceStatus'
];

let missingElements = [];
modalElements.forEach(elementId => {
    const element = document.getElementById(elementId);
    if (!element) {
        missingElements.push(elementId);
    }
});

if (missingElements.length === 0) {
    console.log('✅ All modal form elements found');
} else {
    console.log('❌ Missing modal elements:', missingElements);
}

// Test 6: Test opening platform edit modal
try {
    console.log('Testing modal opening...');
    scenarioBuilder.openPlatformEditModal(testPlatform);

    // Check if modal is visible
    const modal = document.getElementById('platformEditModal');
    if (modal && modal.style.display === 'block') {
        console.log('✅ Platform edit modal opened successfully');

        // Test closing modal
        scenarioBuilder.closePlatformEditModal();
        if (modal.style.display === 'none') {
            console.log('✅ Platform edit modal closed successfully');
        } else {
            console.log('❌ Modal did not close properly');
        }
    } else {
        console.log('❌ Modal did not open properly');
    }
} catch (error) {
    console.log('❌ Error testing modal functionality:', error);
}

// Test 7: Test notification system
try {
    scenarioBuilder.showNotification('Test notification - success', 'success');
    console.log('✅ Success notification shown');

    setTimeout(() => {
        scenarioBuilder.showNotification('Test notification - info', 'info');
        console.log('✅ Info notification shown');
    }, 1000);

    setTimeout(() => {
        scenarioBuilder.showNotification('Test notification - error', 'error');
        console.log('✅ Error notification shown');
    }, 2000);
} catch (error) {
    console.log('❌ Error testing notification system:', error);
}

console.log('=== Test Complete ===');
console.log('If all tests passed, the enhanced platform editing modal system is working correctly!');
