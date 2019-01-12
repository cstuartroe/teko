import os

for root, dirs, files in os.walk(".", topdown=False):
    for filename in files:
        filepath = os.path.join(root,filename)
        if filename.endswith(".cpp"):
            with open(filepath,"r") as fh:
                dirty = fh.read()
            num_tabs = dirty.count("\t")
            if num_tabs > 0:
                clean = dirty.replace("\t"," "*4)
                with open(filepath,"w") as fh:
                    fh.write(clean)
                print("Replaced",num_tabs,"tabs in",filepath)
