default:
	echo "make what? aws, gcp, or all?"
	exit 2

all:
	$(MAKE) aws
	$(MAKE) gcp

aws:
	cd aws &&$(MAKE) all

gcp:
	cd gcp &&$(MAKE) all

clean:
	cd aws &&$(MAKE) clean
	cd gcp &&$(MAKE) clean
