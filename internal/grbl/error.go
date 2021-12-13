package grbl

import "fmt"

const (
	ErrorExpectedCommandLetter       = 1
	ErrorBadNumberFormat             = 2
	ErrorInvalidStatement            = 3
	ErrorNegativeValue               = 4
	ErrorSettingDisabled             = 5
	ErrorSettingStepPulseMin         = 6
	ErrorSettingReadFail             = 7
	ErrorIdleError                   = 8
	ErrorSystemGcLock                = 9
	ErrorSoftLimitError              = 10
	ErrorOverflow                    = 11
	ErrorMaxStepRateExceeded         = 12
	ErrorCheckDoor                   = 13
	ErrorLineLengthExceeded          = 14
	ErrorTravelExceeded              = 15
	ErrorInvalidJogCommand           = 16
	ErrorSettingDisabledLaser        = 17
	ErrorGcodeUnsupportedCommand     = 20
	ErrorGcodeModalGroupViolation    = 21
	ErrorGcodeUndefinedFeedRate      = 22
	ErrorGcodeCommandValueNotInteger = 23
	ErrorGcodeAxisCommandConflict    = 24
	ErrorGcodeWordRepeated           = 25
	ErrorGcodeNoAxisWords            = 26
	ErrorGcodeInvalidLineNumber      = 27
	ErrorGcodeValueWordMissing       = 28
	ErrorGcodeUnsupportedCoordSys    = 29
	ErrorGcodeG53InvalidMotionMode   = 30
	ErrorGcodeAxisWordsExist         = 31
	ErrorGcodeNoAxisWordsInPlane     = 32
	ErrorGcodeInvalidTarget          = 33
	ErrorGcodeArcRadiusError         = 34
	ErrorGcodeNoOffsetsInPlane       = 35
	ErrorGcodeUnusedWords            = 36
	ErrorGcodeG43DynamicAxisError    = 37
	ErrorGcodeMaxValueExceeded       = 38
)

var (
	errorMap = map[uint8]string{
		ErrorExpectedCommandLetter:       "G-code words consist of a letter and a value. Letter was not found",
		ErrorBadNumberFormat:             "Missing the expected G-code word value or numeric value format is not valid",
		ErrorInvalidStatement:            "Grbl '$' system command was not recognized or supported",
		ErrorNegativeValue:               "Negative value received for an expected positive value",
		ErrorSettingDisabled:             "Homing cycle failure. Homing is not enabled via settings",
		ErrorSettingStepPulseMin:         "Minimum step pulse time must be greater than 3usec",
		ErrorSettingReadFail:             "An EEPROM read failed. Auto-restoring affected EEPROM to default values",
		ErrorIdleError:                   "Grbl '$' command cannot be used unless Grbl is IDLE. Ensures smooth operation during a job",
		ErrorSystemGcLock:                "G-code commands are locked out during alarm or jog state",
		ErrorSoftLimitError:              "Soft limits cannot be enabled without homing also enabled",
		ErrorOverflow:                    "Max characters per line exceeded. Received command line was not executed",
		ErrorMaxStepRateExceeded:         "Grbl '$' setting value cause the step rate to exceed the maximum supported",
		ErrorCheckDoor:                   "Safety door detected as opened and door state initiated",
		ErrorLineLengthExceeded:          "Build info or startup line exceeded EEPROM line length limit. Line not stored",
		ErrorTravelExceeded:              "Jog target exceeds machine travel. Jog command has been ignored",
		ErrorInvalidJogCommand:           "Jog command has no '=' or contains prohibited g-code",
		ErrorSettingDisabledLaser:        "Laser mode requires PWM output",
		ErrorGcodeUnsupportedCommand:     "Unsupported or invalid g-code command found in block",
		ErrorGcodeModalGroupViolation:    "More than one g-code command from same modal group found in block",
		ErrorGcodeUndefinedFeedRate:      "Feed rate has not yet been set or is undefined",
		ErrorGcodeCommandValueNotInteger: "G-code command in block requires an integer value",
		ErrorGcodeAxisCommandConflict:    "More than one g-code command that requires axis words found in block",
		ErrorGcodeWordRepeated:           "Repeated g-code word found in block",
		ErrorGcodeNoAxisWords:            "No axis words found in block for g-code command or current modal state which requires them",
		ErrorGcodeInvalidLineNumber:      "Line number value is invalid",
		ErrorGcodeValueWordMissing:       "G-code command is missing a required value word",
		ErrorGcodeUnsupportedCoordSys:    "G59.x work coordinate systems are not supported",
		ErrorGcodeG53InvalidMotionMode:   "G53 only allowed with G0 and G1 motion modes",
		ErrorGcodeAxisWordsExist:         "Axis words found in block when no command or current modal state uses them",
		ErrorGcodeNoAxisWordsInPlane:     "G2 and G3 arcs require at least one in-plane axis word",
		ErrorGcodeInvalidTarget:          "Motion command target is invalid",
		ErrorGcodeArcRadiusError:         "Arc radius value is invalid",
		ErrorGcodeNoOffsetsInPlane:       "G2 and G3 arcs require at least one in-plane offset word",
		ErrorGcodeUnusedWords:            "Unused value words found in block",
		ErrorGcodeG43DynamicAxisError:    "G43.1 dynamic tool length offset is not assigned to configured tool length axis",
		ErrorGcodeMaxValueExceeded:       "Tool number greater than max supported value",
	}
)

type Error uint8

func NewError(id uint8) Error {
	return 0
}

func (e Error) Error() string {
	if msg, found := errorMap[uint8(e)]; found {
		return "grbl: error: " + msg
	}
	return fmt.Sprintf("grbl: error: unknown (%d)", e)
}
