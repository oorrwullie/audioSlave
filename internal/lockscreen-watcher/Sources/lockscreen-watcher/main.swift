import Foundation

let center = DistributedNotificationCenter.default()

func handle(_ message: String) {
    print(message)
    fflush(stdout)
}

center.addObserver(forName: NSNotification.Name("com.apple.screenIsLocked"), object: nil, queue: nil) { _ in
    handle("LOCKED")
}

center.addObserver(forName: NSNotification.Name("com.apple.screenIsUnlocked"), object: nil, queue: nil) { _ in
    handle("UNLOCKED")
}

RunLoop.main.run()
