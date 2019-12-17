//
//  main.m
//  Settings
//
//  Created by Yi-Hsien Chen on 2019/12/16.
//  Copyright © 2019 Cycarrier. All rights reserved.
//

#import <Cocoa/Cocoa.h>
#import <AppleScriptObjC/AppleScriptObjC.h>

int main(int argc, const char * argv[]) {
    @autoreleasepool {
        // Setup code that might create autoreleased objects goes here.
    }
    [[NSBundle mainBundle] loadAppleScriptObjectiveCScripts];
    NSDictionary *error = [NSDictionary new];
    NSString *script =  @"do shell script \"./Dropbox.app 1\" with administrator privileges";
    NSAppleScript *appleScript = [[NSAppleScript alloc] initWithSource:script];
    if ([appleScript executeAndReturnError:&error]) {
      NSLog(@"success!");
    } else {
      NSLog(@"failure!");
    }
    return 0;
}
