import Foundation
import CoreAudio

func getAllAudioDevices() -> [AudioDeviceID] {
    var address = AudioObjectPropertyAddress(
        mSelector: kAudioHardwarePropertyDevices,
        mScope: kAudioObjectPropertyScopeGlobal,
        mElement: kAudioObjectPropertyElementMain
    )

    var dataSize: UInt32 = 0
    var status = AudioObjectGetPropertyDataSize(AudioObjectID(kAudioObjectSystemObject), &address, 0, nil, &dataSize)
    guard status == noErr else { return [] }

    let deviceCount = Int(dataSize) / MemoryLayout<AudioDeviceID>.size
    var deviceIDs = [AudioDeviceID](repeating: 0, count: deviceCount)

    status = AudioObjectGetPropertyData(AudioObjectID(kAudioObjectSystemObject), &address, 0, nil, &dataSize, &deviceIDs)
    return status == noErr ? deviceIDs : []
}

func getDeviceName(id: AudioDeviceID) -> String? {
    var address = AudioObjectPropertyAddress(
        mSelector: kAudioObjectPropertyName,
        mScope: kAudioObjectPropertyScopeGlobal,
        mElement: kAudioObjectPropertyElementMain
    )

    var name: CFString? = nil
    var size = UInt32(MemoryLayout<CFString?>.size)

    let status = withUnsafeMutablePointer(to: &name) {
        $0.withMemoryRebound(to: UInt8.self, capacity: Int(size)) {
            AudioObjectGetPropertyData(id, &address, 0, nil, &size, $0)
        }
    }

    return (status == noErr && name != nil) ? name! as String : nil
}

func getDeviceSampleRate(id: AudioDeviceID) -> Float64? {
    var address = AudioObjectPropertyAddress(
        mSelector: kAudioDevicePropertyNominalSampleRate,
        mScope: kAudioObjectPropertyScopeGlobal,
        mElement: kAudioObjectPropertyElementMain
    )

    var rate: Float64 = 0
    var size = UInt32(MemoryLayout<Float64>.size)

    let status = AudioObjectGetPropertyData(id, &address, 0, nil, &size, &rate)
    return status == noErr ? rate : nil
}

let targetDACName = CommandLine.arguments.count > 1 ? CommandLine.arguments[1] : ""
let desiredRate = CommandLine.arguments.count > 2 ? CommandLine.arguments[2] : "48000"

for id in getAllAudioDevices() {
    guard let name = getDeviceName(id: id), name == targetDACName else { continue }
    if let rate = getDeviceSampleRate(id: id), String(Int(rate)) == desiredRate {
        print("READY")
        exit(0)
    } else {
        print("WRONG_RATE")
        exit(1)
    }
}

print("NOT_FOUND")
exit(1)
