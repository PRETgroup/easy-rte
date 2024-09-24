import subprocess
import os 
build_targets = [
    ["windows", "amd64"],
    ["darwin", "amd64"],
    ["darwin", "arm64"],
    ["linux", "amd64"],
    ["linux", "arm"],
    ["linux", "arm64"],
]

my_env = os.environ.copy()

#run go get github.com/PRETgroup/goFB/goFB
subprocess.run(["go", "get", "github.com/PRETgroup/goFB/goFB"], env=my_env)
  
# Create easy-rte-c and easy-rte-parser executables for each platform
for target in build_targets:
    my_env["GOOS"] = target[0]
    my_env["GOARCH"] = target[1]
    print(f"Building for {target[0]} {target[1]}")
    output_filename_c = "easy-rte-c"
    output_filename_parser = "easy-rte-parser"
    if target[0] == "windows":
        output_filename_c += ".exe"
        output_filename_parser += ".exe"
    #create a directory for each targets
    os.makedirs(f"bin/{target[0]}_{target[1]}", exist_ok=True)
    #build easy-rte-c from ./rtec/main e.g.
    # go build -o easy-rte-c ./rtec/main
    
    subprocess.run(["go", "build", "-o", f"bin/{target[0]}_{target[1]}/{output_filename_c}", "./rtec/main"], env=my_env)
    #build easy-rte-parser from ./rteparser/main e.g.
    # go build -o easy-rte-parser ./rteparser/main
    subprocess.run(["go", "build", "-o", f"bin/{target[0]}_{target[1]}/{output_filename_parser}", "./rteparser/main"], env=my_env)
    # zip the target directory
    # zip -r bin/linux_amd64.zip bin/linux_amd64 from the bin directory
    subprocess.run(["zip", "-r", f"{target[0]}_{target[1]}.zip", f"{target[0]}_{target[1]}"], env=my_env, cwd="bin")
    
# done
print("Done!")