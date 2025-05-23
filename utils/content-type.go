package utils

// DB - Mime is a collection of mime types with extension as key and content-type as value.
var DB = map[string]struct {
	ContentType  string
	Compressible bool
}{
	"123": {
		ContentType:  "application/vnd.lotus-1-2-3",
		Compressible: false,
	},
	"3dml": {
		ContentType:  "text/vnd.in3d.3dml",
		Compressible: false,
	},
	"3ds": {
		ContentType:  "image/x-3ds",
		Compressible: false,
	},
	"3g2": {
		ContentType:  "video/3gpp2",
		Compressible: false,
	},
	"3gp": {
		ContentType:  "video/3gpp",
		Compressible: false,
	},
	"3gpp": {
		ContentType:  "video/3gpp",
		Compressible: false,
	},
	"7z": {
		ContentType:  "application/x-7z-compressed",
		Compressible: false,
	},
	"aab": {
		ContentType:  "application/x-authorware-bin",
		Compressible: false,
	},
	"aac": {
		ContentType:  "audio/x-aac",
		Compressible: false,
	},
	"aam": {
		ContentType:  "application/x-authorware-map",
		Compressible: false,
	},
	"aas": {
		ContentType:  "application/x-authorware-seg",
		Compressible: false,
	},
	"abw": {
		ContentType:  "application/x-abiword",
		Compressible: false,
	},
	"ac": {
		ContentType:  "application/pkix-attr-cert",
		Compressible: false,
	},
	"acc": {
		ContentType:  "application/vnd.americandynamics.acc",
		Compressible: false,
	},
	"ace": {
		ContentType:  "application/x-ace-compressed",
		Compressible: false,
	},
	"acu": {
		ContentType:  "application/vnd.acucobol",
		Compressible: false,
	},
	"acutc": {
		ContentType:  "application/vnd.acucorp",
		Compressible: false,
	},
	"adp": {
		ContentType:  "audio/adpcm",
		Compressible: false,
	},
	"aep": {
		ContentType:  "application/vnd.audiograph",
		Compressible: false,
	},
	"afm": {
		ContentType:  "application/x-font-type1",
		Compressible: false,
	},
	"afp": {
		ContentType:  "application/vnd.ibm.modcap",
		Compressible: false,
	},
	"ahead": {
		ContentType:  "application/vnd.ahead.space",
		Compressible: false,
	},
	"ai": {
		ContentType:  "application/postscript",
		Compressible: false,
	},
	"aif": {
		ContentType:  "audio/x-aiff",
		Compressible: false,
	},
	"aifc": {
		ContentType:  "audio/x-aiff",
		Compressible: false,
	},
	"aiff": {
		ContentType:  "audio/x-aiff",
		Compressible: false,
	},
	"air": {
		ContentType:  "application/vnd.adobe.air-application-installer-package+zip",
		Compressible: false,
	},
	"ait": {
		ContentType:  "application/vnd.dvb.ait",
		Compressible: false,
	},
	"ami": {
		ContentType:  "application/vnd.amiga.ami",
		Compressible: false,
	},
	"apk": {
		ContentType:  "application/vnd.android.package-archive",
		Compressible: false,
	},
	"apng": {
		ContentType:  "image/apng",
		Compressible: false,
	},
	"appcache": {
		ContentType:  "text/cache-manifest",
		Compressible: false,
	},
	"application": {
		ContentType:  "application/x-ms-application",
		Compressible: false,
	},
	"apr": {
		ContentType:  "application/vnd.lotus-approach",
		Compressible: false,
	},
	"arc": {
		ContentType:  "application/x-freearc",
		Compressible: false,
	},
	"arj": {
		ContentType:  "application/x-arj",
		Compressible: false,
	},
	"asc": {
		ContentType:  "application/pgp-signature",
		Compressible: false,
	},
	"asf": {
		ContentType:  "video/x-ms-asf",
		Compressible: false,
	},
	"asm": {
		ContentType:  "text/x-asm",
		Compressible: false,
	},
	"aso": {
		ContentType:  "application/vnd.accpac.simply.aso",
		Compressible: false,
	},
	"asx": {
		ContentType:  "video/x-ms-asf",
		Compressible: false,
	},
	"atc": {
		ContentType:  "application/vnd.acucorp",
		Compressible: false,
	},
	"atom": {
		ContentType:  "application/atom+xml",
		Compressible: false,
	},
	"atomcat": {
		ContentType:  "application/atomcat+xml",
		Compressible: false,
	},
	"atomsvc": {
		ContentType:  "application/atomsvc+xml",
		Compressible: false,
	},
	"atx": {
		ContentType:  "application/vnd.antix.game-component",
		Compressible: false,
	},
	"au": {
		ContentType:  "audio/basic",
		Compressible: false,
	},
	"avi": {
		ContentType:  "video/x-msvideo",
		Compressible: false,
	},
	"aw": {
		ContentType:  "application/applixware",
		Compressible: false,
	},
	"azf": {
		ContentType:  "application/vnd.airzip.filesecure.azf",
		Compressible: false,
	},
	"azs": {
		ContentType:  "application/vnd.airzip.filesecure.azs",
		Compressible: false,
	},
	"azv": {
		ContentType:  "image/vnd.airzip.accelerator.azv",
		Compressible: false,
	},
	"azw": {
		ContentType:  "application/vnd.amazon.ebook",
		Compressible: false,
	},
	"bat": {
		ContentType:  "application/x-msdownload",
		Compressible: false,
	},
	"bcpio": {
		ContentType:  "application/x-bcpio",
		Compressible: false,
	},
	"bdf": {
		ContentType:  "application/x-font-bdf",
		Compressible: false,
	},
	"bdm": {
		ContentType:  "application/vnd.syncml.dm+wbxml",
		Compressible: false,
	},
	"bdoc": {
		ContentType:  "application/x-bdoc",
		Compressible: false,
	},
	"bed": {
		ContentType:  "application/vnd.realvnc.bed",
		Compressible: false,
	},
	"bh2": {
		ContentType:  "application/vnd.fujitsu.oasysprs",
		Compressible: false,
	},
	"bin": {
		ContentType:  "application/octet-stream",
		Compressible: false,
	},
	"blb": {
		ContentType:  "application/x-blorb",
		Compressible: false,
	},
	"blorb": {
		ContentType:  "application/x-blorb",
		Compressible: false,
	},
	"bmi": {
		ContentType:  "application/vnd.bmi",
		Compressible: false,
	},
	"bmp": {
		ContentType:  "image/x-ms-bmp",
		Compressible: false,
	},
	"book": {
		ContentType:  "application/vnd.framemaker",
		Compressible: false,
	},
	"box": {
		ContentType:  "application/vnd.previewsystems.box",
		Compressible: false,
	},
	"boz": {
		ContentType:  "application/x-bzip2",
		Compressible: false,
	},
	"bpk": {
		ContentType:  "application/octet-stream",
		Compressible: false,
	},
	"btif": {
		ContentType:  "image/prs.btif",
		Compressible: false,
	},
	"buffer": {
		ContentType:  "application/octet-stream",
		Compressible: false,
	},
	"bz": {
		ContentType:  "application/x-bzip",
		Compressible: false,
	},
	"bz2": {
		ContentType:  "application/x-bzip2",
		Compressible: false,
	},
	"c": {
		ContentType:  "text/x-c",
		Compressible: false,
	},
	"c11amc": {
		ContentType:  "application/vnd.cluetrust.cartomobile-config",
		Compressible: false,
	},
	"c11amz": {
		ContentType:  "application/vnd.cluetrust.cartomobile-config-pkg",
		Compressible: false,
	},
	"c4d": {
		ContentType:  "application/vnd.clonk.c4group",
		Compressible: false,
	},
	"c4f": {
		ContentType:  "application/vnd.clonk.c4group",
		Compressible: false,
	},
	"c4g": {
		ContentType:  "application/vnd.clonk.c4group",
		Compressible: false,
	},
	"c4p": {
		ContentType:  "application/vnd.clonk.c4group",
		Compressible: false,
	},
	"c4u": {
		ContentType:  "application/vnd.clonk.c4group",
		Compressible: false,
	},
	"cab": {
		ContentType:  "application/vnd.ms-cab-compressed",
		Compressible: false,
	},
	"caf": {
		ContentType:  "audio/x-caf",
		Compressible: false,
	},
	"cap": {
		ContentType:  "application/vnd.tcpdump.pcap",
		Compressible: false,
	},
	"car": {
		ContentType:  "application/vnd.curl.car",
		Compressible: false,
	},
	"cat": {
		ContentType:  "application/vnd.ms-pki.seccat",
		Compressible: false,
	},
	"cb7": {
		ContentType:  "application/x-cbr",
		Compressible: false,
	},
	"cba": {
		ContentType:  "application/x-cbr",
		Compressible: false,
	},
	"cbr": {
		ContentType:  "application/x-cbr",
		Compressible: false,
	},
	"cbt": {
		ContentType:  "application/x-cbr",
		Compressible: false,
	},
	"cbz": {
		ContentType:  "application/x-cbr",
		Compressible: false,
	},
	"cc": {
		ContentType:  "text/x-c",
		Compressible: false,
	},
	"cco": {
		ContentType:  "application/x-cocoa",
		Compressible: false,
	},
	"cct": {
		ContentType:  "application/x-director",
		Compressible: false,
	},
	"ccxml": {
		ContentType:  "application/ccxml+xml",
		Compressible: false,
	},
	"cdbcmsg": {
		ContentType:  "application/vnd.contact.cmsg",
		Compressible: false,
	},
	"cdf": {
		ContentType:  "application/x-netcdf",
		Compressible: false,
	},
	"cdkey": {
		ContentType:  "application/vnd.mediastation.cdkey",
		Compressible: false,
	},
	"cdmia": {
		ContentType:  "application/cdmi-capability",
		Compressible: false,
	},
	"cdmic": {
		ContentType:  "application/cdmi-container",
		Compressible: false,
	},
	"cdmid": {
		ContentType:  "application/cdmi-domain",
		Compressible: false,
	},
	"cdmio": {
		ContentType:  "application/cdmi-object",
		Compressible: false,
	},
	"cdmiq": {
		ContentType:  "application/cdmi-queue",
		Compressible: false,
	},
	"cdx": {
		ContentType:  "chemical/x-cdx",
		Compressible: false,
	},
	"cdxml": {
		ContentType:  "application/vnd.chemdraw+xml",
		Compressible: false,
	},
	"cdy": {
		ContentType:  "application/vnd.cinderella",
		Compressible: false,
	},
	"cer": {
		ContentType:  "application/pkix-cert",
		Compressible: false,
	},
	"cfs": {
		ContentType:  "application/x-cfs-compressed",
		Compressible: false,
	},
	"cgm": {
		ContentType:  "image/cgm",
		Compressible: false,
	},
	"chat": {
		ContentType:  "application/x-chat",
		Compressible: false,
	},
	"chm": {
		ContentType:  "application/vnd.ms-htmlhelp",
		Compressible: false,
	},
	"chrt": {
		ContentType:  "application/vnd.kde.kchart",
		Compressible: false,
	},
	"cif": {
		ContentType:  "chemical/x-cif",
		Compressible: false,
	},
	"cii": {
		ContentType:  "application/vnd.anser-web-certificate-issue-initiation",
		Compressible: false,
	},
	"cil": {
		ContentType:  "application/vnd.ms-artgalry",
		Compressible: false,
	},
	"cla": {
		ContentType:  "application/vnd.claymore",
		Compressible: false,
	},
	"class": {
		ContentType:  "application/java-vm",
		Compressible: false,
	},
	"clkk": {
		ContentType:  "application/vnd.crick.clicker.keyboard",
		Compressible: false,
	},
	"clkp": {
		ContentType:  "application/vnd.crick.clicker.palette",
		Compressible: false,
	},
	"clkt": {
		ContentType:  "application/vnd.crick.clicker.template",
		Compressible: false,
	},
	"clkw": {
		ContentType:  "application/vnd.crick.clicker.wordbank",
		Compressible: false,
	},
	"clkx": {
		ContentType:  "application/vnd.crick.clicker",
		Compressible: false,
	},
	"clp": {
		ContentType:  "application/x-msclip",
		Compressible: false,
	},
	"cmc": {
		ContentType:  "application/vnd.cosmocaller",
		Compressible: false,
	},
	"cmdf": {
		ContentType:  "chemical/x-cmdf",
		Compressible: false,
	},
	"cml": {
		ContentType:  "chemical/x-cml",
		Compressible: false,
	},
	"cmp": {
		ContentType:  "application/vnd.yellowriver-custom-menu",
		Compressible: false,
	},
	"cmx": {
		ContentType:  "image/x-cmx",
		Compressible: false,
	},
	"cod": {
		ContentType:  "application/vnd.rim.cod",
		Compressible: false,
	},
	"coffee": {
		ContentType:  "text/coffeescript",
		Compressible: false,
	},
	"com": {
		ContentType:  "application/x-msdownload",
		Compressible: false,
	},
	"conf": {
		ContentType:  "text/plain",
		Compressible: false,
	},
	"cpio": {
		ContentType:  "application/x-cpio",
		Compressible: false,
	},
	"cpp": {
		ContentType:  "text/x-c",
		Compressible: false,
	},
	"cpt": {
		ContentType:  "application/mac-compactpro",
		Compressible: false,
	},
	"crd": {
		ContentType:  "application/x-mscardfile",
		Compressible: false,
	},
	"crl": {
		ContentType:  "application/pkix-crl",
		Compressible: false,
	},
	"crt": {
		ContentType:  "application/x-x509-ca-cert",
		Compressible: false,
	},
	"crx": {
		ContentType:  "application/x-chrome-extension",
		Compressible: false,
	},
	"cryptonote": {
		ContentType:  "application/vnd.rig.cryptonote",
		Compressible: false,
	},
	"csh": {
		ContentType:  "application/x-csh",
		Compressible: false,
	},
	"csl": {
		ContentType:  "application/vnd.citationstyles.style+xml",
		Compressible: false,
	},
	"csml": {
		ContentType:  "chemical/x-csml",
		Compressible: false,
	},
	"csp": {
		ContentType:  "application/vnd.commonspace",
		Compressible: false,
	},
	"css": {
		ContentType:  "text/css",
		Compressible: false,
	},
	"cst": {
		ContentType:  "application/x-director",
		Compressible: false,
	},
	"csv": {
		ContentType:  "text/csv",
		Compressible: false,
	},
	"cu": {
		ContentType:  "application/cu-seeme",
		Compressible: false,
	},
	"curl": {
		ContentType:  "text/vnd.curl",
		Compressible: false,
	},
	"cww": {
		ContentType:  "application/prs.cww",
		Compressible: false,
	},
	"cxt": {
		ContentType:  "application/x-director",
		Compressible: false,
	},
	"cxx": {
		ContentType:  "text/x-c",
		Compressible: false,
	},
	"dae": {
		ContentType:  "model/vnd.collada+xml",
		Compressible: false,
	},
	"daf": {
		ContentType:  "application/vnd.mobius.daf",
		Compressible: false,
	},
	"dart": {
		ContentType:  "application/vnd.dart",
		Compressible: false,
	},
	"dataless": {
		ContentType:  "application/vnd.fdsn.seed",
		Compressible: false,
	},
	"davmount": {
		ContentType:  "application/davmount+xml",
		Compressible: false,
	},
	"dbk": {
		ContentType:  "application/docbook+xml",
		Compressible: false,
	},
	"dcr": {
		ContentType:  "application/x-director",
		Compressible: false,
	},
	"dcurl": {
		ContentType:  "text/vnd.curl.dcurl",
		Compressible: false,
	},
	"dd2": {
		ContentType:  "application/vnd.oma.dd2+xml",
		Compressible: false,
	},
	"ddd": {
		ContentType:  "application/vnd.fujixerox.ddd",
		Compressible: false,
	},
	"deb": {
		ContentType:  "application/x-debian-package",
		Compressible: false,
	},
	"def": {
		ContentType:  "text/plain",
		Compressible: false,
	},
	"deploy": {
		ContentType:  "application/octet-stream",
		Compressible: false,
	},
	"der": {
		ContentType:  "application/x-x509-ca-cert",
		Compressible: false,
	},
	"dfac": {
		ContentType:  "application/vnd.dreamfactory",
		Compressible: false,
	},
	"dgc": {
		ContentType:  "application/x-dgc-compressed",
		Compressible: false,
	},
	"dic": {
		ContentType:  "text/x-c",
		Compressible: false,
	},
	"dir": {
		ContentType:  "application/x-director",
		Compressible: false,
	},
	"dis": {
		ContentType:  "application/vnd.mobius.dis",
		Compressible: false,
	},
	"disposition-notification": {
		ContentType:  "message/disposition-notification",
		Compressible: false,
	},
	"dist": {
		ContentType:  "application/octet-stream",
		Compressible: false,
	},
	"distz": {
		ContentType:  "application/octet-stream",
		Compressible: false,
	},
	"djv": {
		ContentType:  "image/vnd.djvu",
		Compressible: false,
	},
	"djvu": {
		ContentType:  "image/vnd.djvu",
		Compressible: false,
	},
	"dll": {
		ContentType:  "application/x-msdownload",
		Compressible: false,
	},
	"dmg": {
		ContentType:  "application/x-apple-diskimage",
		Compressible: false,
	},
	"dmp": {
		ContentType:  "application/vnd.tcpdump.pcap",
		Compressible: false,
	},
	"dms": {
		ContentType:  "application/octet-stream",
		Compressible: false,
	},
	"dna": {
		ContentType:  "application/vnd.dna",
		Compressible: false,
	},
	"doc": {
		ContentType:  "application/msword",
		Compressible: false,
	},
	"docm": {
		ContentType:  "application/vnd.ms-word.document.macroenabled.12",
		Compressible: false,
	},
	"docx": {
		ContentType:  "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
		Compressible: false,
	},
	"dot": {
		ContentType:  "application/msword",
		Compressible: false,
	},
	"dotm": {
		ContentType:  "application/vnd.ms-word.template.macroenabled.12",
		Compressible: false,
	},
	"dotx": {
		ContentType:  "application/vnd.openxmlformats-officedocument.wordprocessingml.template",
		Compressible: false,
	},
	"dp": {
		ContentType:  "application/vnd.osgi.dp",
		Compressible: false,
	},
	"dpg": {
		ContentType:  "application/vnd.dpgraph",
		Compressible: false,
	},
	"dra": {
		ContentType:  "audio/vnd.dra",
		Compressible: false,
	},
	"drle": {
		ContentType:  "image/dicom-rle",
		Compressible: false,
	},
	"dsc": {
		ContentType:  "text/prs.lines.tag",
		Compressible: false,
	},
	"dssc": {
		ContentType:  "application/dssc+der",
		Compressible: false,
	},
	"dtb": {
		ContentType:  "application/x-dtbook+xml",
		Compressible: false,
	},
	"dtd": {
		ContentType:  "application/xml-dtd",
		Compressible: false,
	},
	"dts": {
		ContentType:  "audio/vnd.dts",
		Compressible: false,
	},
	"dtshd": {
		ContentType:  "audio/vnd.dts.hd",
		Compressible: false,
	},
	"dump": {
		ContentType:  "application/octet-stream",
		Compressible: false,
	},
	"dvb": {
		ContentType:  "video/vnd.dvb.file",
		Compressible: false,
	},
	"dvi": {
		ContentType:  "application/x-dvi",
		Compressible: false,
	},
	"dwf": {
		ContentType:  "model/vnd.dwf",
		Compressible: false,
	},
	"dwg": {
		ContentType:  "image/vnd.dwg",
		Compressible: false,
	},
	"dxf": {
		ContentType:  "image/vnd.dxf",
		Compressible: false,
	},
	"dxp": {
		ContentType:  "application/vnd.spotfire.dxp",
		Compressible: false,
	},
	"dxr": {
		ContentType:  "application/x-director",
		Compressible: false,
	},
	"ear": {
		ContentType:  "application/java-archive",
		Compressible: false,
	},
	"ecelp4800": {
		ContentType:  "audio/vnd.nuera.ecelp4800",
		Compressible: false,
	},
	"ecelp7470": {
		ContentType:  "audio/vnd.nuera.ecelp7470",
		Compressible: false,
	},
	"ecelp9600": {
		ContentType:  "audio/vnd.nuera.ecelp9600",
		Compressible: false,
	},
	"ecma": {
		ContentType:  "application/ecmascript",
		Compressible: false,
	},
	"edm": {
		ContentType:  "application/vnd.novadigm.edm",
		Compressible: false,
	},
	"edx": {
		ContentType:  "application/vnd.novadigm.edx",
		Compressible: false,
	},
	"efif": {
		ContentType:  "application/vnd.picsel",
		Compressible: false,
	},
	"ei6": {
		ContentType:  "application/vnd.pg.osasli",
		Compressible: false,
	},
	"elc": {
		ContentType:  "application/octet-stream",
		Compressible: false,
	},
	"emf": {
		ContentType:  "image/emf",
		Compressible: false,
	},
	"eml": {
		ContentType:  "message/rfc822",
		Compressible: false,
	},
	"emma": {
		ContentType:  "application/emma+xml",
		Compressible: false,
	},
	"emz": {
		ContentType:  "application/x-msmetafile",
		Compressible: false,
	},
	"eol": {
		ContentType:  "audio/vnd.digital-winds",
		Compressible: false,
	},
	"eot": {
		ContentType:  "application/vnd.ms-fontobject",
		Compressible: false,
	},
	"eps": {
		ContentType:  "application/postscript",
		Compressible: false,
	},
	"epub": {
		ContentType:  "application/epub+zip",
		Compressible: false,
	},
	"es": {
		ContentType:  "application/ecmascript",
		Compressible: false,
	},
	"es3": {
		ContentType:  "application/vnd.eszigno3+xml",
		Compressible: false,
	},
	"esa": {
		ContentType:  "application/vnd.osgi.subsystem",
		Compressible: false,
	},
	"esf": {
		ContentType:  "application/vnd.epson.esf",
		Compressible: false,
	},
	"et3": {
		ContentType:  "application/vnd.eszigno3+xml",
		Compressible: false,
	},
	"etx": {
		ContentType:  "text/x-setext",
		Compressible: false,
	},
	"eva": {
		ContentType:  "application/x-eva",
		Compressible: false,
	},
	"evy": {
		ContentType:  "application/x-envoy",
		Compressible: false,
	},
	"exe": {
		ContentType:  "application/x-msdownload",
		Compressible: false,
	},
	"exi": {
		ContentType:  "application/exi",
		Compressible: false,
	},
	"exr": {
		ContentType:  "image/aces",
		Compressible: false,
	},
	"ext": {
		ContentType:  "application/vnd.novadigm.ext",
		Compressible: false,
	},
	"ez": {
		ContentType:  "application/andrew-inset",
		Compressible: false,
	},
	"ez2": {
		ContentType:  "application/vnd.ezpix-album",
		Compressible: false,
	},
	"ez3": {
		ContentType:  "application/vnd.ezpix-package",
		Compressible: false,
	},
	"f": {
		ContentType:  "text/x-fortran",
		Compressible: false,
	},
	"f4v": {
		ContentType:  "video/x-f4v",
		Compressible: false,
	},
	"f77": {
		ContentType:  "text/x-fortran",
		Compressible: false,
	},
	"f90": {
		ContentType:  "text/x-fortran",
		Compressible: false,
	},
	"fbs": {
		ContentType:  "image/vnd.fastbidsheet",
		Compressible: false,
	},
	"fcdt": {
		ContentType:  "application/vnd.adobe.formscentral.fcdt",
		Compressible: false,
	},
	"fcs": {
		ContentType:  "application/vnd.isac.fcs",
		Compressible: false,
	},
	"fdf": {
		ContentType:  "application/vnd.fdf",
		Compressible: false,
	},
	"fe_launch": {
		ContentType:  "application/vnd.denovo.fcselayout-link",
		Compressible: false,
	},
	"fg5": {
		ContentType:  "application/vnd.fujitsu.oasysgp",
		Compressible: false,
	},
	"fgd": {
		ContentType:  "application/x-director",
		Compressible: false,
	},
	"fh": {
		ContentType:  "image/x-freehand",
		Compressible: false,
	},
	"fh4": {
		ContentType:  "image/x-freehand",
		Compressible: false,
	},
	"fh5": {
		ContentType:  "image/x-freehand",
		Compressible: false,
	},
	"fh7": {
		ContentType:  "image/x-freehand",
		Compressible: false,
	},
	"fhc": {
		ContentType:  "image/x-freehand",
		Compressible: false,
	},
	"fig": {
		ContentType:  "application/x-xfig",
		Compressible: false,
	},
	"fits": {
		ContentType:  "image/fits",
		Compressible: false,
	},
	"flac": {
		ContentType:  "audio/x-flac",
		Compressible: false,
	},
	"fli": {
		ContentType:  "video/x-fli",
		Compressible: false,
	},
	"flo": {
		ContentType:  "application/vnd.micrografx.flo",
		Compressible: false,
	},
	"flv": {
		ContentType:  "video/x-flv",
		Compressible: false,
	},
	"flw": {
		ContentType:  "application/vnd.kde.kivio",
		Compressible: false,
	},
	"flx": {
		ContentType:  "text/vnd.fmi.flexstor",
		Compressible: false,
	},
	"fly": {
		ContentType:  "text/vnd.fly",
		Compressible: false,
	},
	"fm": {
		ContentType:  "application/vnd.framemaker",
		Compressible: false,
	},
	"fnc": {
		ContentType:  "application/vnd.frogans.fnc",
		Compressible: false,
	},
	"for": {
		ContentType:  "text/x-fortran",
		Compressible: false,
	},
	"fpx": {
		ContentType:  "image/vnd.fpx",
		Compressible: false,
	},
	"frame": {
		ContentType:  "application/vnd.framemaker",
		Compressible: false,
	},
	"fsc": {
		ContentType:  "application/vnd.fsc.weblaunch",
		Compressible: false,
	},
	"fst": {
		ContentType:  "image/vnd.fst",
		Compressible: false,
	},
	"ftc": {
		ContentType:  "application/vnd.fluxtime.clip",
		Compressible: false,
	},
	"fti": {
		ContentType:  "application/vnd.anser-web-funds-transfer-initiation",
		Compressible: false,
	},
	"fvt": {
		ContentType:  "video/vnd.fvt",
		Compressible: false,
	},
	"fxp": {
		ContentType:  "application/vnd.adobe.fxp",
		Compressible: false,
	},
	"fxpl": {
		ContentType:  "application/vnd.adobe.fxp",
		Compressible: false,
	},
	"fzs": {
		ContentType:  "application/vnd.fuzzysheet",
		Compressible: false,
	},
	"g2w": {
		ContentType:  "application/vnd.geoplan",
		Compressible: false,
	},
	"g3": {
		ContentType:  "image/g3fax",
		Compressible: false,
	},
	"g3w": {
		ContentType:  "application/vnd.geospace",
		Compressible: false,
	},
	"gac": {
		ContentType:  "application/vnd.groove-account",
		Compressible: false,
	},
	"gam": {
		ContentType:  "application/x-tads",
		Compressible: false,
	},
	"gbr": {
		ContentType:  "application/rpki-ghostbusters",
		Compressible: false,
	},
	"gca": {
		ContentType:  "application/x-gca-compressed",
		Compressible: false,
	},
	"gdl": {
		ContentType:  "model/vnd.gdl",
		Compressible: false,
	},
	"gdoc": {
		ContentType:  "application/vnd.google-apps.document",
		Compressible: false,
	},
	"geo": {
		ContentType:  "application/vnd.dynageo",
		Compressible: false,
	},
	"geojson": {
		ContentType:  "application/geo+json",
		Compressible: false,
	},
	"gex": {
		ContentType:  "application/vnd.geometry-explorer",
		Compressible: false,
	},
	"ggb": {
		ContentType:  "application/vnd.geogebra.file",
		Compressible: false,
	},
	"ggt": {
		ContentType:  "application/vnd.geogebra.tool",
		Compressible: false,
	},
	"ghf": {
		ContentType:  "application/vnd.groove-help",
		Compressible: false,
	},
	"gif": {
		ContentType:  "image/gif",
		Compressible: false,
	},
	"gim": {
		ContentType:  "application/vnd.groove-identity-message",
		Compressible: false,
	},
	"glb": {
		ContentType:  "model/gltf-binary",
		Compressible: false,
	},
	"gltf": {
		ContentType:  "model/gltf+json",
		Compressible: false,
	},
	"gml": {
		ContentType:  "application/gml+xml",
		Compressible: false,
	},
	"gmx": {
		ContentType:  "application/vnd.gmx",
		Compressible: false,
	},
	"gnumeric": {
		ContentType:  "application/x-gnumeric",
		Compressible: false,
	},
	"gph": {
		ContentType:  "application/vnd.flographit",
		Compressible: false,
	},
	"gpx": {
		ContentType:  "application/gpx+xml",
		Compressible: false,
	},
	"gqf": {
		ContentType:  "application/vnd.grafeq",
		Compressible: false,
	},
	"gqs": {
		ContentType:  "application/vnd.grafeq",
		Compressible: false,
	},
	"gram": {
		ContentType:  "application/srgs",
		Compressible: false,
	},
	"gramps": {
		ContentType:  "application/x-gramps-xml",
		Compressible: false,
	},
	"gre": {
		ContentType:  "application/vnd.geometry-explorer",
		Compressible: false,
	},
	"grv": {
		ContentType:  "application/vnd.groove-injector",
		Compressible: false,
	},
	"grxml": {
		ContentType:  "application/srgs+xml",
		Compressible: false,
	},
	"gsf": {
		ContentType:  "application/x-font-ghostscript",
		Compressible: false,
	},
	"gsheet": {
		ContentType:  "application/vnd.google-apps.spreadsheet",
		Compressible: false,
	},
	"gslides": {
		ContentType:  "application/vnd.google-apps.presentation",
		Compressible: false,
	},
	"gtar": {
		ContentType:  "application/x-gtar",
		Compressible: false,
	},
	"gtm": {
		ContentType:  "application/vnd.groove-tool-message",
		Compressible: false,
	},
	"gtw": {
		ContentType:  "model/vnd.gtw",
		Compressible: false,
	},
	"gv": {
		ContentType:  "text/vnd.graphviz",
		Compressible: false,
	},
	"gxf": {
		ContentType:  "application/gxf",
		Compressible: false,
	},
	"gxt": {
		ContentType:  "application/vnd.geonext",
		Compressible: false,
	},
	"gz": {
		ContentType:  "application/gzip",
		Compressible: false,
	},
	"h": {
		ContentType:  "text/x-c",
		Compressible: false,
	},
	"h261": {
		ContentType:  "video/h261",
		Compressible: false,
	},
	"h263": {
		ContentType:  "video/h263",
		Compressible: false,
	},
	"h264": {
		ContentType:  "video/h264",
		Compressible: false,
	},
	"hal": {
		ContentType:  "application/vnd.hal+xml",
		Compressible: false,
	},
	"hbci": {
		ContentType:  "application/vnd.hbci",
		Compressible: false,
	},
	"hbs": {
		ContentType:  "text/x-handlebars-template",
		Compressible: false,
	},
	"hdd": {
		ContentType:  "application/x-virtualbox-hdd",
		Compressible: false,
	},
	"hdf": {
		ContentType:  "application/x-hdf",
		Compressible: false,
	},
	"heic": {
		ContentType:  "image/heic",
		Compressible: false,
	},
	"heics": {
		ContentType:  "image/heic-sequence",
		Compressible: false,
	},
	"heif": {
		ContentType:  "image/heif",
		Compressible: false,
	},
	"heifs": {
		ContentType:  "image/heif-sequence",
		Compressible: false,
	},
	"hh": {
		ContentType:  "text/x-c",
		Compressible: false,
	},
	"hjson": {
		ContentType:  "application/hjson",
		Compressible: false,
	},
	"hlp": {
		ContentType:  "application/winhlp",
		Compressible: false,
	},
	"hpgl": {
		ContentType:  "application/vnd.hp-hpgl",
		Compressible: false,
	},
	"hpid": {
		ContentType:  "application/vnd.hp-hpid",
		Compressible: false,
	},
	"hps": {
		ContentType:  "application/vnd.hp-hps",
		Compressible: false,
	},
	"hqx": {
		ContentType:  "application/mac-binhex40",
		Compressible: false,
	},
	"htc": {
		ContentType:  "text/x-component",
		Compressible: false,
	},
	"htke": {
		ContentType:  "application/vnd.kenameaapp",
		Compressible: false,
	},
	"htm": {
		ContentType:  "text/html",
		Compressible: false,
	},
	"html": {
		ContentType:  "text/html",
		Compressible: false,
	},
	"hvd": {
		ContentType:  "application/vnd.yamaha.hv-dic",
		Compressible: false,
	},
	"hvp": {
		ContentType:  "application/vnd.yamaha.hv-voice",
		Compressible: false,
	},
	"hvs": {
		ContentType:  "application/vnd.yamaha.hv-script",
		Compressible: false,
	},
	"i2g": {
		ContentType:  "application/vnd.intergeo",
		Compressible: false,
	},
	"icc": {
		ContentType:  "application/vnd.iccprofile",
		Compressible: false,
	},
	"ice": {
		ContentType:  "x-conference/x-cooltalk",
		Compressible: false,
	},
	"icm": {
		ContentType:  "application/vnd.iccprofile",
		Compressible: false,
	},
	"ico": {
		ContentType:  "image/x-icon",
		Compressible: false,
	},
	"ics": {
		ContentType:  "text/calendar",
		Compressible: false,
	},
	"ief": {
		ContentType:  "image/ief",
		Compressible: false,
	},
	"ifb": {
		ContentType:  "text/calendar",
		Compressible: false,
	},
	"ifm": {
		ContentType:  "application/vnd.shana.informed.formdata",
		Compressible: false,
	},
	"iges": {
		ContentType:  "model/iges",
		Compressible: false,
	},
	"igl": {
		ContentType:  "application/vnd.igloader",
		Compressible: false,
	},
	"igm": {
		ContentType:  "application/vnd.insors.igm",
		Compressible: false,
	},
	"igs": {
		ContentType:  "model/iges",
		Compressible: false,
	},
	"igx": {
		ContentType:  "application/vnd.micrografx.igx",
		Compressible: false,
	},
	"iif": {
		ContentType:  "application/vnd.shana.informed.interchange",
		Compressible: false,
	},
	"img": {
		ContentType:  "application/octet-stream",
		Compressible: false,
	},
	"imp": {
		ContentType:  "application/vnd.accpac.simply.imp",
		Compressible: false,
	},
	"ims": {
		ContentType:  "application/vnd.ms-ims",
		Compressible: false,
	},
	"in": {
		ContentType:  "text/plain",
		Compressible: false,
	},
	"ini": {
		ContentType:  "text/plain",
		Compressible: false,
	},
	"ink": {
		ContentType:  "application/inkml+xml",
		Compressible: false,
	},
	"inkml": {
		ContentType:  "application/inkml+xml",
		Compressible: false,
	},
	"install": {
		ContentType:  "application/x-install-instructions",
		Compressible: false,
	},
	"iota": {
		ContentType:  "application/vnd.astraea-software.iota",
		Compressible: false,
	},
	"ipfix": {
		ContentType:  "application/ipfix",
		Compressible: false,
	},
	"ipk": {
		ContentType:  "application/vnd.shana.informed.package",
		Compressible: false,
	},
	"irm": {
		ContentType:  "application/vnd.ibm.rights-management",
		Compressible: false,
	},
	"irp": {
		ContentType:  "application/vnd.irepository.package+xml",
		Compressible: false,
	},
	"iso": {
		ContentType:  "application/x-iso9660-image",
		Compressible: false,
	},
	"itp": {
		ContentType:  "application/vnd.shana.informed.formtemplate",
		Compressible: false,
	},
	"ivp": {
		ContentType:  "application/vnd.immervision-ivp",
		Compressible: false,
	},
	"ivu": {
		ContentType:  "application/vnd.immervision-ivu",
		Compressible: false,
	},
	"jad": {
		ContentType:  "text/vnd.sun.j2me.app-descriptor",
		Compressible: false,
	},
	"jade": {
		ContentType:  "text/jade",
		Compressible: false,
	},
	"jam": {
		ContentType:  "application/vnd.jam",
		Compressible: false,
	},
	"jar": {
		ContentType:  "application/java-archive",
		Compressible: false,
	},
	"jardiff": {
		ContentType:  "application/x-java-archive-diff",
		Compressible: false,
	},
	"java": {
		ContentType:  "text/x-java-source",
		Compressible: false,
	},
	"jisp": {
		ContentType:  "application/vnd.jisp",
		Compressible: false,
	},
	"jls": {
		ContentType:  "image/jls",
		Compressible: false,
	},
	"jlt": {
		ContentType:  "application/vnd.hp-jlyt",
		Compressible: false,
	},
	"jng": {
		ContentType:  "image/x-jng",
		Compressible: false,
	},
	"jnlp": {
		ContentType:  "application/x-java-jnlp-file",
		Compressible: false,
	},
	"joda": {
		ContentType:  "application/vnd.joost.joda-archive",
		Compressible: false,
	},
	"jp2": {
		ContentType:  "image/jp2",
		Compressible: false,
	},
	"jpe": {
		ContentType:  "image/jpeg",
		Compressible: false,
	},
	"jpeg": {
		ContentType:  "image/jpeg",
		Compressible: false,
	},
	"jpf": {
		ContentType:  "image/jpx",
		Compressible: false,
	},
	"jpg": {
		ContentType:  "image/jpeg",
		Compressible: false,
	},
	"jpg2": {
		ContentType:  "image/jp2",
		Compressible: false,
	},
	"jpgm": {
		ContentType:  "video/jpm",
		Compressible: false,
	},
	"jpgv": {
		ContentType:  "video/jpeg",
		Compressible: false,
	},
	"jpm": {
		ContentType:  "video/jpm",
		Compressible: false,
	},
	"jpx": {
		ContentType:  "image/jpx",
		Compressible: false,
	},
	"js": {
		ContentType:  "application/javascript",
		Compressible: false,
	},
	"json": {
		ContentType:  "application/json",
		Compressible: false,
	},
	"json5": {
		ContentType:  "application/json5",
		Compressible: false,
	},
	"jsonld": {
		ContentType:  "application/ld+json",
		Compressible: false,
	},
	"jsonml": {
		ContentType:  "application/jsonml+json",
		Compressible: false,
	},
	"jsx": {
		ContentType:  "text/jsx",
		Compressible: false,
	},
	"kar": {
		ContentType:  "audio/midi",
		Compressible: false,
	},
	"karbon": {
		ContentType:  "application/vnd.kde.karbon",
		Compressible: false,
	},
	"keynote": {
		ContentType:  "application/vnd.apple.keynote",
		Compressible: false,
	},
	"kfo": {
		ContentType:  "application/vnd.kde.kformula",
		Compressible: false,
	},
	"kia": {
		ContentType:  "application/vnd.kidspiration",
		Compressible: false,
	},
	"kml": {
		ContentType:  "application/vnd.google-earth.kml+xml",
		Compressible: false,
	},
	"kmz": {
		ContentType:  "application/vnd.google-earth.kmz",
		Compressible: false,
	},
	"kne": {
		ContentType:  "application/vnd.kinar",
		Compressible: false,
	},
	"knp": {
		ContentType:  "application/vnd.kinar",
		Compressible: false,
	},
	"kon": {
		ContentType:  "application/vnd.kde.kontour",
		Compressible: false,
	},
	"kpr": {
		ContentType:  "application/vnd.kde.kpresenter",
		Compressible: false,
	},
	"kpt": {
		ContentType:  "application/vnd.kde.kpresenter",
		Compressible: false,
	},
	"kpxx": {
		ContentType:  "application/vnd.ds-keypoint",
		Compressible: false,
	},
	"ksp": {
		ContentType:  "application/vnd.kde.kspread",
		Compressible: false,
	},
	"ktr": {
		ContentType:  "application/vnd.kahootz",
		Compressible: false,
	},
	"ktx": {
		ContentType:  "image/ktx",
		Compressible: false,
	},
	"ktz": {
		ContentType:  "application/vnd.kahootz",
		Compressible: false,
	},
	"kwd": {
		ContentType:  "application/vnd.kde.kword",
		Compressible: false,
	},
	"kwt": {
		ContentType:  "application/vnd.kde.kword",
		Compressible: false,
	},
	"lasxml": {
		ContentType:  "application/vnd.las.las+xml",
		Compressible: false,
	},
	"latex": {
		ContentType:  "application/x-latex",
		Compressible: false,
	},
	"lbd": {
		ContentType:  "application/vnd.llamagraphics.life-balance.desktop",
		Compressible: false,
	},
	"lbe": {
		ContentType:  "application/vnd.llamagraphics.life-balance.exchange+xml",
		Compressible: false,
	},
	"les": {
		ContentType:  "application/vnd.hhe.lesson-player",
		Compressible: false,
	},
	"less": {
		ContentType:  "text/less",
		Compressible: false,
	},
	"lha": {
		ContentType:  "application/x-lzh-compressed",
		Compressible: false,
	},
	"link66": {
		ContentType:  "application/vnd.route66.link66+xml",
		Compressible: false,
	},
	"list": {
		ContentType:  "text/plain",
		Compressible: false,
	},
	"list3820": {
		ContentType:  "application/vnd.ibm.modcap",
		Compressible: false,
	},
	"listafp": {
		ContentType:  "application/vnd.ibm.modcap",
		Compressible: false,
	},
	"litcoffee": {
		ContentType:  "text/coffeescript",
		Compressible: false,
	},
	"lnk": {
		ContentType:  "application/x-ms-shortcut",
		Compressible: false,
	},
	"log": {
		ContentType:  "text/plain",
		Compressible: false,
	},
	"lostxml": {
		ContentType:  "application/lost+xml",
		Compressible: false,
	},
	"lrf": {
		ContentType:  "application/octet-stream",
		Compressible: false,
	},
	"lrm": {
		ContentType:  "application/vnd.ms-lrm",
		Compressible: false,
	},
	"ltf": {
		ContentType:  "application/vnd.frogans.ltf",
		Compressible: false,
	},
	"lua": {
		ContentType:  "text/x-lua",
		Compressible: false,
	},
	"luac": {
		ContentType:  "application/x-lua-bytecode",
		Compressible: false,
	},
	"lvp": {
		ContentType:  "audio/vnd.lucent.voice",
		Compressible: false,
	},
	"lwp": {
		ContentType:  "application/vnd.lotus-wordpro",
		Compressible: false,
	},
	"lzh": {
		ContentType:  "application/x-lzh-compressed",
		Compressible: false,
	},
	"m13": {
		ContentType:  "application/x-msmediaview",
		Compressible: false,
	},
	"m14": {
		ContentType:  "application/x-msmediaview",
		Compressible: false,
	},
	"m1v": {
		ContentType:  "video/mpeg",
		Compressible: false,
	},
	"m21": {
		ContentType:  "application/mp21",
		Compressible: false,
	},
	"m2a": {
		ContentType:  "audio/mpeg",
		Compressible: false,
	},
	"m2v": {
		ContentType:  "video/mpeg",
		Compressible: false,
	},
	"m3a": {
		ContentType:  "audio/mpeg",
		Compressible: false,
	},
	"m3u": {
		ContentType:  "audio/x-mpegurl",
		Compressible: false,
	},
	"m3u8": {
		ContentType:  "application/vnd.apple.mpegurl",
		Compressible: false,
	},
	"m4a": {
		ContentType:  "audio/x-m4a",
		Compressible: false,
	},
	"m4p": {
		ContentType:  "application/mp4",
		Compressible: false,
	},
	"m4u": {
		ContentType:  "video/vnd.mpegurl",
		Compressible: false,
	},
	"m4v": {
		ContentType:  "video/x-m4v",
		Compressible: false,
	},
	"ma": {
		ContentType:  "application/mathematica",
		Compressible: false,
	},
	"mads": {
		ContentType:  "application/mads+xml",
		Compressible: false,
	},
	"mag": {
		ContentType:  "application/vnd.ecowin.chart",
		Compressible: false,
	},
	"maker": {
		ContentType:  "application/vnd.framemaker",
		Compressible: false,
	},
	"man": {
		ContentType:  "text/troff",
		Compressible: false,
	},
	"manifest": {
		ContentType:  "text/cache-manifest",
		Compressible: false,
	},
	"map": {
		ContentType:  "application/json",
		Compressible: false,
	},
	"mar": {
		ContentType:  "application/octet-stream",
		Compressible: false,
	},
	"markdown": {
		ContentType:  "text/markdown",
		Compressible: false,
	},
	"mathml": {
		ContentType:  "application/mathml+xml",
		Compressible: false,
	},
	"mb": {
		ContentType:  "application/mathematica",
		Compressible: false,
	},
	"mbk": {
		ContentType:  "application/vnd.mobius.mbk",
		Compressible: false,
	},
	"mbox": {
		ContentType:  "application/mbox",
		Compressible: false,
	},
	"mc1": {
		ContentType:  "application/vnd.medcalcdata",
		Compressible: false,
	},
	"mcd": {
		ContentType:  "application/vnd.mcd",
		Compressible: false,
	},
	"mcurl": {
		ContentType:  "text/vnd.curl.mcurl",
		Compressible: false,
	},
	"md": {
		ContentType:  "text/markdown",
		Compressible: false,
	},
	"mdb": {
		ContentType:  "application/x-msaccess",
		Compressible: false,
	},
	"mdi": {
		ContentType:  "image/vnd.ms-modi",
		Compressible: false,
	},
	"me": {
		ContentType:  "text/troff",
		Compressible: false,
	},
	"mesh": {
		ContentType:  "model/mesh",
		Compressible: false,
	},
	"meta4": {
		ContentType:  "application/metalink4+xml",
		Compressible: false,
	},
	"metalink": {
		ContentType:  "application/metalink+xml",
		Compressible: false,
	},
	"mets": {
		ContentType:  "application/mets+xml",
		Compressible: false,
	},
	"mfm": {
		ContentType:  "application/vnd.mfmp",
		Compressible: false,
	},
	"mft": {
		ContentType:  "application/rpki-manifest",
		Compressible: false,
	},
	"mgp": {
		ContentType:  "application/vnd.osgeo.mapguide.package",
		Compressible: false,
	},
	"mgz": {
		ContentType:  "application/vnd.proteus.magazine",
		Compressible: false,
	},
	"mid": {
		ContentType:  "audio/midi",
		Compressible: false,
	},
	"midi": {
		ContentType:  "audio/midi",
		Compressible: false,
	},
	"mie": {
		ContentType:  "application/x-mie",
		Compressible: false,
	},
	"mif": {
		ContentType:  "application/vnd.mif",
		Compressible: false,
	},
	"mime": {
		ContentType:  "message/rfc822",
		Compressible: false,
	},
	"mj2": {
		ContentType:  "video/mj2",
		Compressible: false,
	},
	"mjp2": {
		ContentType:  "video/mj2",
		Compressible: false,
	},
	"mjs": {
		ContentType:  "application/javascript",
		Compressible: false,
	},
	"mk3d": {
		ContentType:  "video/x-matroska",
		Compressible: false,
	},
	"mka": {
		ContentType:  "audio/x-matroska",
		Compressible: false,
	},
	"mkd": {
		ContentType:  "text/x-markdown",
		Compressible: false,
	},
	"mks": {
		ContentType:  "video/x-matroska",
		Compressible: false,
	},
	"mkv": {
		ContentType:  "video/x-matroska",
		Compressible: false,
	},
	"mlp": {
		ContentType:  "application/vnd.dolby.mlp",
		Compressible: false,
	},
	"mmd": {
		ContentType:  "application/vnd.chipnuts.karaoke-mmd",
		Compressible: false,
	},
	"mmf": {
		ContentType:  "application/vnd.smaf",
		Compressible: false,
	},
	"mml": {
		ContentType:  "text/mathml",
		Compressible: false,
	},
	"mmr": {
		ContentType:  "image/vnd.fujixerox.edmics-mmr",
		Compressible: false,
	},
	"mng": {
		ContentType:  "video/x-mng",
		Compressible: false,
	},
	"mny": {
		ContentType:  "application/x-msmoney",
		Compressible: false,
	},
	"mobi": {
		ContentType:  "application/x-mobipocket-ebook",
		Compressible: false,
	},
	"mods": {
		ContentType:  "application/mods+xml",
		Compressible: false,
	},
	"mov": {
		ContentType:  "video/quicktime",
		Compressible: false,
	},
	"movie": {
		ContentType:  "video/x-sgi-movie",
		Compressible: false,
	},
	"mp2": {
		ContentType:  "audio/mpeg",
		Compressible: false,
	},
	"mp21": {
		ContentType:  "application/mp21",
		Compressible: false,
	},
	"mp2a": {
		ContentType:  "audio/mpeg",
		Compressible: false,
	},
	"mp3": {
		ContentType:  "audio/mpeg",
		Compressible: false,
	},
	"mp4": {
		ContentType:  "video/mp4",
		Compressible: false,
	},
	"mp4a": {
		ContentType:  "audio/mp4",
		Compressible: false,
	},
	"mp4s": {
		ContentType:  "application/mp4",
		Compressible: false,
	},
	"mp4v": {
		ContentType:  "video/mp4",
		Compressible: false,
	},
	"mpc": {
		ContentType:  "application/vnd.mophun.certificate",
		Compressible: false,
	},
	"mpd": {
		ContentType:  "application/dash+xml",
		Compressible: false,
	},
	"mpe": {
		ContentType:  "video/mpeg",
		Compressible: false,
	},
	"mpeg": {
		ContentType:  "video/mpeg",
		Compressible: false,
	},
	"mpg": {
		ContentType:  "video/mpeg",
		Compressible: false,
	},
	"mpg4": {
		ContentType:  "video/mp4",
		Compressible: false,
	},
	"mpga": {
		ContentType:  "audio/mpeg",
		Compressible: false,
	},
	"mpkg": {
		ContentType:  "application/vnd.apple.installer+xml",
		Compressible: false,
	},
	"mpm": {
		ContentType:  "application/vnd.blueice.multipass",
		Compressible: false,
	},
	"mpn": {
		ContentType:  "application/vnd.mophun.application",
		Compressible: false,
	},
	"mpp": {
		ContentType:  "application/vnd.ms-project",
		Compressible: false,
	},
	"mpt": {
		ContentType:  "application/vnd.ms-project",
		Compressible: false,
	},
	"mpy": {
		ContentType:  "application/vnd.ibm.minipay",
		Compressible: false,
	},
	"mqy": {
		ContentType:  "application/vnd.mobius.mqy",
		Compressible: false,
	},
	"mrc": {
		ContentType:  "application/marc",
		Compressible: false,
	},
	"mrcx": {
		ContentType:  "application/marcxml+xml",
		Compressible: false,
	},
	"ms": {
		ContentType:  "text/troff",
		Compressible: false,
	},
	"mscml": {
		ContentType:  "application/mediaservercontrol+xml",
		Compressible: false,
	},
	"mseed": {
		ContentType:  "application/vnd.fdsn.mseed",
		Compressible: false,
	},
	"mseq": {
		ContentType:  "application/vnd.mseq",
		Compressible: false,
	},
	"msf": {
		ContentType:  "application/vnd.epson.msf",
		Compressible: false,
	},
	"msg": {
		ContentType:  "application/vnd.ms-outlook",
		Compressible: false,
	},
	"msh": {
		ContentType:  "model/mesh",
		Compressible: false,
	},
	"msi": {
		ContentType:  "application/x-msdownload",
		Compressible: false,
	},
	"msl": {
		ContentType:  "application/vnd.mobius.msl",
		Compressible: false,
	},
	"msm": {
		ContentType:  "application/octet-stream",
		Compressible: false,
	},
	"msp": {
		ContentType:  "application/octet-stream",
		Compressible: false,
	},
	"msty": {
		ContentType:  "application/vnd.muvee.style",
		Compressible: false,
	},
	"mts": {
		ContentType:  "model/vnd.mts",
		Compressible: false,
	},
	"mus": {
		ContentType:  "application/vnd.musician",
		Compressible: false,
	},
	"musicxml": {
		ContentType:  "application/vnd.recordare.musicxml+xml",
		Compressible: false,
	},
	"mvb": {
		ContentType:  "application/x-msmediaview",
		Compressible: false,
	},
	"mwf": {
		ContentType:  "application/vnd.mfer",
		Compressible: false,
	},
	"mxf": {
		ContentType:  "application/mxf",
		Compressible: false,
	},
	"mxl": {
		ContentType:  "application/vnd.recordare.musicxml",
		Compressible: false,
	},
	"mxml": {
		ContentType:  "application/xv+xml",
		Compressible: false,
	},
	"mxs": {
		ContentType:  "application/vnd.triscape.mxs",
		Compressible: false,
	},
	"mxu": {
		ContentType:  "video/vnd.mpegurl",
		Compressible: false,
	},
	"n-gage": {
		ContentType:  "application/vnd.nokia.n-gage.symbian.install",
		Compressible: false,
	},
	"n3": {
		ContentType:  "text/n3",
		Compressible: false,
	},
	"nb": {
		ContentType:  "application/mathematica",
		Compressible: false,
	},
	"nbp": {
		ContentType:  "application/vnd.wolfram.player",
		Compressible: false,
	},
	"nc": {
		ContentType:  "application/x-netcdf",
		Compressible: false,
	},
	"ncx": {
		ContentType:  "application/x-dtbncx+xml",
		Compressible: false,
	},
	"nfo": {
		ContentType:  "text/x-nfo",
		Compressible: false,
	},
	"ngdat": {
		ContentType:  "application/vnd.nokia.n-gage.data",
		Compressible: false,
	},
	"nitf": {
		ContentType:  "application/vnd.nitf",
		Compressible: false,
	},
	"nlu": {
		ContentType:  "application/vnd.neurolanguage.nlu",
		Compressible: false,
	},
	"nml": {
		ContentType:  "application/vnd.enliven",
		Compressible: false,
	},
	"nnd": {
		ContentType:  "application/vnd.noblenet-directory",
		Compressible: false,
	},
	"nns": {
		ContentType:  "application/vnd.noblenet-sealer",
		Compressible: false,
	},
	"nnw": {
		ContentType:  "application/vnd.noblenet-web",
		Compressible: false,
	},
	"npx": {
		ContentType:  "image/vnd.net-fpx",
		Compressible: false,
	},
	"nsc": {
		ContentType:  "application/x-conference",
		Compressible: false,
	},
	"nsf": {
		ContentType:  "application/vnd.lotus-notes",
		Compressible: false,
	},
	"ntf": {
		ContentType:  "application/vnd.nitf",
		Compressible: false,
	},
	"numbers": {
		ContentType:  "application/vnd.apple.numbers",
		Compressible: false,
	},
	"nzb": {
		ContentType:  "application/x-nzb",
		Compressible: false,
	},
	"oa2": {
		ContentType:  "application/vnd.fujitsu.oasys2",
		Compressible: false,
	},
	"oa3": {
		ContentType:  "application/vnd.fujitsu.oasys3",
		Compressible: false,
	},
	"oas": {
		ContentType:  "application/vnd.fujitsu.oasys",
		Compressible: false,
	},
	"obd": {
		ContentType:  "application/x-msbinder",
		Compressible: false,
	},
	"obj": {
		ContentType:  "application/x-tgif",
		Compressible: false,
	},
	"oda": {
		ContentType:  "application/oda",
		Compressible: false,
	},
	"odb": {
		ContentType:  "application/vnd.oasis.opendocument.database",
		Compressible: false,
	},
	"odc": {
		ContentType:  "application/vnd.oasis.opendocument.chart",
		Compressible: false,
	},
	"odf": {
		ContentType:  "application/vnd.oasis.opendocument.formula",
		Compressible: false,
	},
	"odft": {
		ContentType:  "application/vnd.oasis.opendocument.formula-template",
		Compressible: false,
	},
	"odg": {
		ContentType:  "application/vnd.oasis.opendocument.graphics",
		Compressible: false,
	},
	"odi": {
		ContentType:  "application/vnd.oasis.opendocument.image",
		Compressible: false,
	},
	"odm": {
		ContentType:  "application/vnd.oasis.opendocument.text-master",
		Compressible: false,
	},
	"odp": {
		ContentType:  "application/vnd.oasis.opendocument.presentation",
		Compressible: false,
	},
	"ods": {
		ContentType:  "application/vnd.oasis.opendocument.spreadsheet",
		Compressible: false,
	},
	"odt": {
		ContentType:  "application/vnd.oasis.opendocument.text",
		Compressible: false,
	},
	"oga": {
		ContentType:  "audio/ogg",
		Compressible: false,
	},
	"ogg": {
		ContentType:  "audio/ogg",
		Compressible: false,
	},
	"ogv": {
		ContentType:  "video/ogg",
		Compressible: false,
	},
	"ogx": {
		ContentType:  "application/ogg",
		Compressible: false,
	},
	"omdoc": {
		ContentType:  "application/omdoc+xml",
		Compressible: false,
	},
	"onepkg": {
		ContentType:  "application/onenote",
		Compressible: false,
	},
	"onetmp": {
		ContentType:  "application/onenote",
		Compressible: false,
	},
	"onetoc": {
		ContentType:  "application/onenote",
		Compressible: false,
	},
	"onetoc2": {
		ContentType:  "application/onenote",
		Compressible: false,
	},
	"opf": {
		ContentType:  "application/oebps-package+xml",
		Compressible: false,
	},
	"opml": {
		ContentType:  "text/x-opml",
		Compressible: false,
	},
	"oprc": {
		ContentType:  "application/vnd.palm",
		Compressible: false,
	},
	"org": {
		ContentType:  "text/x-org",
		Compressible: false,
	},
	"osf": {
		ContentType:  "application/vnd.yamaha.openscoreformat",
		Compressible: false,
	},
	"osfpvg": {
		ContentType:  "application/vnd.yamaha.openscoreformat.osfpvg+xml",
		Compressible: false,
	},
	"otc": {
		ContentType:  "application/vnd.oasis.opendocument.chart-template",
		Compressible: false,
	},
	"otf": {
		ContentType:  "font/otf",
		Compressible: false,
	},
	"otg": {
		ContentType:  "application/vnd.oasis.opendocument.graphics-template",
		Compressible: false,
	},
	"oth": {
		ContentType:  "application/vnd.oasis.opendocument.text-web",
		Compressible: false,
	},
	"oti": {
		ContentType:  "application/vnd.oasis.opendocument.image-template",
		Compressible: false,
	},
	"otp": {
		ContentType:  "application/vnd.oasis.opendocument.presentation-template",
		Compressible: false,
	},
	"ots": {
		ContentType:  "application/vnd.oasis.opendocument.spreadsheet-template",
		Compressible: false,
	},
	"ott": {
		ContentType:  "application/vnd.oasis.opendocument.text-template",
		Compressible: false,
	},
	"ova": {
		ContentType:  "application/x-virtualbox-ova",
		Compressible: false,
	},
	"ovf": {
		ContentType:  "application/x-virtualbox-ovf",
		Compressible: false,
	},
	"owl": {
		ContentType:  "application/rdf+xml",
		Compressible: false,
	},
	"oxps": {
		ContentType:  "application/oxps",
		Compressible: false,
	},
	"oxt": {
		ContentType:  "application/vnd.openofficeorg.extension",
		Compressible: false,
	},
	"p": {
		ContentType:  "text/x-pascal",
		Compressible: false,
	},
	"p10": {
		ContentType:  "application/pkcs10",
		Compressible: false,
	},
	"p12": {
		ContentType:  "application/x-pkcs12",
		Compressible: false,
	},
	"p7b": {
		ContentType:  "application/x-pkcs7-certificates",
		Compressible: false,
	},
	"p7c": {
		ContentType:  "application/pkcs7-mime",
		Compressible: false,
	},
	"p7m": {
		ContentType:  "application/pkcs7-mime",
		Compressible: false,
	},
	"p7r": {
		ContentType:  "application/x-pkcs7-certreqresp",
		Compressible: false,
	},
	"p7s": {
		ContentType:  "application/pkcs7-signature",
		Compressible: false,
	},
	"p8": {
		ContentType:  "application/pkcs8",
		Compressible: false,
	},
	"pac": {
		ContentType:  "application/x-ns-proxy-autoconfig",
		Compressible: false,
	},
	"pages": {
		ContentType:  "application/vnd.apple.pages",
		Compressible: false,
	},
	"pas": {
		ContentType:  "text/x-pascal",
		Compressible: false,
	},
	"paw": {
		ContentType:  "application/vnd.pawaafile",
		Compressible: false,
	},
	"pbd": {
		ContentType:  "application/vnd.powerbuilder6",
		Compressible: false,
	},
	"pbm": {
		ContentType:  "image/x-portable-bitmap",
		Compressible: false,
	},
	"pcap": {
		ContentType:  "application/vnd.tcpdump.pcap",
		Compressible: false,
	},
	"pcf": {
		ContentType:  "application/x-font-pcf",
		Compressible: false,
	},
	"pcl": {
		ContentType:  "application/vnd.hp-pcl",
		Compressible: false,
	},
	"pclxl": {
		ContentType:  "application/vnd.hp-pclxl",
		Compressible: false,
	},
	"pct": {
		ContentType:  "image/x-pict",
		Compressible: false,
	},
	"pcurl": {
		ContentType:  "application/vnd.curl.pcurl",
		Compressible: false,
	},
	"pcx": {
		ContentType:  "image/x-pcx",
		Compressible: false,
	},
	"pdb": {
		ContentType:  "application/x-pilot",
		Compressible: false,
	},
	"pde": {
		ContentType:  "text/x-processing",
		Compressible: false,
	},
	"pdf": {
		ContentType:  "application/pdf",
		Compressible: false,
	},
	"pem": {
		ContentType:  "application/x-x509-ca-cert",
		Compressible: false,
	},
	"pfa": {
		ContentType:  "application/x-font-type1",
		Compressible: false,
	},
	"pfb": {
		ContentType:  "application/x-font-type1",
		Compressible: false,
	},
	"pfm": {
		ContentType:  "application/x-font-type1",
		Compressible: false,
	},
	"pfr": {
		ContentType:  "application/font-tdpfr",
		Compressible: false,
	},
	"pfx": {
		ContentType:  "application/x-pkcs12",
		Compressible: false,
	},
	"pgm": {
		ContentType:  "image/x-portable-graymap",
		Compressible: false,
	},
	"pgn": {
		ContentType:  "application/x-chess-pgn",
		Compressible: false,
	},
	"pgp": {
		ContentType:  "application/pgp-encrypted",
		Compressible: false,
	},
	"php": {
		ContentType:  "application/x-httpd-php",
		Compressible: false,
	},
	"pic": {
		ContentType:  "image/x-pict",
		Compressible: false,
	},
	"pkg": {
		ContentType:  "application/octet-stream",
		Compressible: false,
	},
	"pki": {
		ContentType:  "application/pkixcmp",
		Compressible: false,
	},
	"pkipath": {
		ContentType:  "application/pkix-pkipath",
		Compressible: false,
	},
	"pkpass": {
		ContentType:  "application/vnd.apple.pkpass",
		Compressible: false,
	},
	"pl": {
		ContentType:  "application/x-perl",
		Compressible: false,
	},
	"plb": {
		ContentType:  "application/vnd.3gpp.pic-bw-large",
		Compressible: false,
	},
	"plc": {
		ContentType:  "application/vnd.mobius.plc",
		Compressible: false,
	},
	"plf": {
		ContentType:  "application/vnd.pocketlearn",
		Compressible: false,
	},
	"pls": {
		ContentType:  "application/pls+xml",
		Compressible: false,
	},
	"pm": {
		ContentType:  "application/x-perl",
		Compressible: false,
	},
	"pml": {
		ContentType:  "application/vnd.ctc-posml",
		Compressible: false,
	},
	"png": {
		ContentType:  "image/png",
		Compressible: false,
	},
	"pnm": {
		ContentType:  "image/x-portable-anymap",
		Compressible: false,
	},
	"portpkg": {
		ContentType:  "application/vnd.macports.portpkg",
		Compressible: false,
	},
	"pot": {
		ContentType:  "application/vnd.ms-powerpoint",
		Compressible: false,
	},
	"potm": {
		ContentType:  "application/vnd.ms-powerpoint.template.macroenabled.12",
		Compressible: false,
	},
	"potx": {
		ContentType:  "application/vnd.openxmlformats-officedocument.presentationml.template",
		Compressible: false,
	},
	"ppam": {
		ContentType:  "application/vnd.ms-powerpoint.addin.macroenabled.12",
		Compressible: false,
	},
	"ppd": {
		ContentType:  "application/vnd.cups-ppd",
		Compressible: false,
	},
	"ppm": {
		ContentType:  "image/x-portable-pixmap",
		Compressible: false,
	},
	"pps": {
		ContentType:  "application/vnd.ms-powerpoint",
		Compressible: false,
	},
	"ppsm": {
		ContentType:  "application/vnd.ms-powerpoint.slideshow.macroenabled.12",
		Compressible: false,
	},
	"ppsx": {
		ContentType:  "application/vnd.openxmlformats-officedocument.presentationml.slideshow",
		Compressible: false,
	},
	"ppt": {
		ContentType:  "application/vnd.ms-powerpoint",
		Compressible: false,
	},
	"pptm": {
		ContentType:  "application/vnd.ms-powerpoint.presentation.macroenabled.12",
		Compressible: false,
	},
	"pptx": {
		ContentType:  "application/vnd.openxmlformats-officedocument.presentationml.presentation",
		Compressible: false,
	},
	"pqa": {
		ContentType:  "application/vnd.palm",
		Compressible: false,
	},
	"prc": {
		ContentType:  "application/x-pilot",
		Compressible: false,
	},
	"pre": {
		ContentType:  "application/vnd.lotus-freelance",
		Compressible: false,
	},
	"prf": {
		ContentType:  "application/pics-rules",
		Compressible: false,
	},
	"ps": {
		ContentType:  "application/postscript",
		Compressible: false,
	},
	"psb": {
		ContentType:  "application/vnd.3gpp.pic-bw-small",
		Compressible: false,
	},
	"psd": {
		ContentType:  "image/vnd.adobe.photoshop",
		Compressible: false,
	},
	"psf": {
		ContentType:  "application/x-font-linux-psf",
		Compressible: false,
	},
	"pskcxml": {
		ContentType:  "application/pskc+xml",
		Compressible: false,
	},
	"pti": {
		ContentType:  "image/prs.pti",
		Compressible: false,
	},
	"ptid": {
		ContentType:  "application/vnd.pvi.ptid1",
		Compressible: false,
	},
	"pub": {
		ContentType:  "application/x-mspublisher",
		Compressible: false,
	},
	"pvb": {
		ContentType:  "application/vnd.3gpp.pic-bw-var",
		Compressible: false,
	},
	"pwn": {
		ContentType:  "application/vnd.3m.post-it-notes",
		Compressible: false,
	},
	"pya": {
		ContentType:  "audio/vnd.ms-playready.media.pya",
		Compressible: false,
	},
	"pyv": {
		ContentType:  "video/vnd.ms-playready.media.pyv",
		Compressible: false,
	},
	"qam": {
		ContentType:  "application/vnd.epson.quickanime",
		Compressible: false,
	},
	"qbo": {
		ContentType:  "application/vnd.intu.qbo",
		Compressible: false,
	},
	"qfx": {
		ContentType:  "application/vnd.intu.qfx",
		Compressible: false,
	},
	"qps": {
		ContentType:  "application/vnd.publishare-delta-tree",
		Compressible: false,
	},
	"qt": {
		ContentType:  "video/quicktime",
		Compressible: false,
	},
	"qwd": {
		ContentType:  "application/vnd.quark.quarkxpress",
		Compressible: false,
	},
	"qwt": {
		ContentType:  "application/vnd.quark.quarkxpress",
		Compressible: false,
	},
	"qxb": {
		ContentType:  "application/vnd.quark.quarkxpress",
		Compressible: false,
	},
	"qxd": {
		ContentType:  "application/vnd.quark.quarkxpress",
		Compressible: false,
	},
	"qxl": {
		ContentType:  "application/vnd.quark.quarkxpress",
		Compressible: false,
	},
	"qxt": {
		ContentType:  "application/vnd.quark.quarkxpress",
		Compressible: false,
	},
	"ra": {
		ContentType:  "audio/x-realaudio",
		Compressible: false,
	},
	"ram": {
		ContentType:  "audio/x-pn-realaudio",
		Compressible: false,
	},
	"raml": {
		ContentType:  "application/raml+yaml",
		Compressible: false,
	},
	"rar": {
		ContentType:  "application/x-rar-compressed",
		Compressible: false,
	},
	"ras": {
		ContentType:  "image/x-cmu-raster",
		Compressible: false,
	},
	"rcprofile": {
		ContentType:  "application/vnd.ipunplugged.rcprofile",
		Compressible: false,
	},
	"rdf": {
		ContentType:  "application/rdf+xml",
		Compressible: false,
	},
	"rdz": {
		ContentType:  "application/vnd.data-vision.rdz",
		Compressible: false,
	},
	"rep": {
		ContentType:  "application/vnd.businessobjects",
		Compressible: false,
	},
	"res": {
		ContentType:  "application/x-dtbresource+xml",
		Compressible: false,
	},
	"rgb": {
		ContentType:  "image/x-rgb",
		Compressible: false,
	},
	"rif": {
		ContentType:  "application/reginfo+xml",
		Compressible: false,
	},
	"rip": {
		ContentType:  "audio/vnd.rip",
		Compressible: false,
	},
	"ris": {
		ContentType:  "application/x-research-info-systems",
		Compressible: false,
	},
	"rl": {
		ContentType:  "application/resource-lists+xml",
		Compressible: false,
	},
	"rlc": {
		ContentType:  "image/vnd.fujixerox.edmics-rlc",
		Compressible: false,
	},
	"rld": {
		ContentType:  "application/resource-lists-diff+xml",
		Compressible: false,
	},
	"rm": {
		ContentType:  "application/vnd.rn-realmedia",
		Compressible: false,
	},
	"rmi": {
		ContentType:  "audio/midi",
		Compressible: false,
	},
	"rmp": {
		ContentType:  "audio/x-pn-realaudio-plugin",
		Compressible: false,
	},
	"rms": {
		ContentType:  "application/vnd.jcp.javame.midlet-rms",
		Compressible: false,
	},
	"rmvb": {
		ContentType:  "application/vnd.rn-realmedia-vbr",
		Compressible: false,
	},
	"rnc": {
		ContentType:  "application/relax-ng-compact-syntax",
		Compressible: false,
	},
	"rng": {
		ContentType:  "application/xml",
		Compressible: false,
	},
	"roa": {
		ContentType:  "application/rpki-roa",
		Compressible: false,
	},
	"roff": {
		ContentType:  "text/troff",
		Compressible: false,
	},
	"rp9": {
		ContentType:  "application/vnd.cloanto.rp9",
		Compressible: false,
	},
	"rpm": {
		ContentType:  "application/x-redhat-package-manager",
		Compressible: false,
	},
	"rpss": {
		ContentType:  "application/vnd.nokia.radio-presets",
		Compressible: false,
	},
	"rpst": {
		ContentType:  "application/vnd.nokia.radio-preset",
		Compressible: false,
	},
	"rq": {
		ContentType:  "application/sparql-query",
		Compressible: false,
	},
	"rs": {
		ContentType:  "application/rls-services+xml",
		Compressible: false,
	},
	"rsd": {
		ContentType:  "application/rsd+xml",
		Compressible: false,
	},
	"rss": {
		ContentType:  "application/rss+xml",
		Compressible: false,
	},
	"rtf": {
		ContentType:  "text/rtf",
		Compressible: false,
	},
	"rtx": {
		ContentType:  "text/richtext",
		Compressible: false,
	},
	"run": {
		ContentType:  "application/x-makeself",
		Compressible: false,
	},
	"s": {
		ContentType:  "text/x-asm",
		Compressible: false,
	},
	"s3m": {
		ContentType:  "audio/s3m",
		Compressible: false,
	},
	"saf": {
		ContentType:  "application/vnd.yamaha.smaf-audio",
		Compressible: false,
	},
	"sass": {
		ContentType:  "text/x-sass",
		Compressible: false,
	},
	"sbml": {
		ContentType:  "application/sbml+xml",
		Compressible: false,
	},
	"sc": {
		ContentType:  "application/vnd.ibm.secure-container",
		Compressible: false,
	},
	"scd": {
		ContentType:  "application/x-msschedule",
		Compressible: false,
	},
	"scm": {
		ContentType:  "application/vnd.lotus-screencam",
		Compressible: false,
	},
	"scq": {
		ContentType:  "application/scvp-cv-request",
		Compressible: false,
	},
	"scs": {
		ContentType:  "application/scvp-cv-response",
		Compressible: false,
	},
	"scss": {
		ContentType:  "text/x-scss",
		Compressible: false,
	},
	"scurl": {
		ContentType:  "text/vnd.curl.scurl",
		Compressible: false,
	},
	"sda": {
		ContentType:  "application/vnd.stardivision.draw",
		Compressible: false,
	},
	"sdc": {
		ContentType:  "application/vnd.stardivision.calc",
		Compressible: false,
	},
	"sdd": {
		ContentType:  "application/vnd.stardivision.impress",
		Compressible: false,
	},
	"sdkd": {
		ContentType:  "application/vnd.solent.sdkm+xml",
		Compressible: false,
	},
	"sdkm": {
		ContentType:  "application/vnd.solent.sdkm+xml",
		Compressible: false,
	},
	"sdp": {
		ContentType:  "application/sdp",
		Compressible: false,
	},
	"sdw": {
		ContentType:  "application/vnd.stardivision.writer",
		Compressible: false,
	},
	"sea": {
		ContentType:  "application/x-sea",
		Compressible: false,
	},
	"see": {
		ContentType:  "application/vnd.seemail",
		Compressible: false,
	},
	"seed": {
		ContentType:  "application/vnd.fdsn.seed",
		Compressible: false,
	},
	"sema": {
		ContentType:  "application/vnd.sema",
		Compressible: false,
	},
	"semd": {
		ContentType:  "application/vnd.semd",
		Compressible: false,
	},
	"semf": {
		ContentType:  "application/vnd.semf",
		Compressible: false,
	},
	"ser": {
		ContentType:  "application/java-serialized-object",
		Compressible: false,
	},
	"setpay": {
		ContentType:  "application/set-payment-initiation",
		Compressible: false,
	},
	"setreg": {
		ContentType:  "application/set-registration-initiation",
		Compressible: false,
	},
	"sfd-hdstx": {
		ContentType:  "application/vnd.hydrostatix.sof-data",
		Compressible: false,
	},
	"sfs": {
		ContentType:  "application/vnd.spotfire.sfs",
		Compressible: false,
	},
	"sfv": {
		ContentType:  "text/x-sfv",
		Compressible: false,
	},
	"sgi": {
		ContentType:  "image/sgi",
		Compressible: false,
	},
	"sgl": {
		ContentType:  "application/vnd.stardivision.writer-global",
		Compressible: false,
	},
	"sgm": {
		ContentType:  "text/sgml",
		Compressible: false,
	},
	"sgml": {
		ContentType:  "text/sgml",
		Compressible: false,
	},
	"sh": {
		ContentType:  "application/x-sh",
		Compressible: false,
	},
	"shar": {
		ContentType:  "application/x-shar",
		Compressible: false,
	},
	"shex": {
		ContentType:  "text/shex",
		Compressible: false,
	},
	"shf": {
		ContentType:  "application/shf+xml",
		Compressible: false,
	},
	"shtml": {
		ContentType:  "text/html",
		Compressible: false,
	},
	"sid": {
		ContentType:  "image/x-mrsid-image",
		Compressible: false,
	},
	"sig": {
		ContentType:  "application/pgp-signature",
		Compressible: false,
	},
	"sil": {
		ContentType:  "audio/silk",
		Compressible: false,
	},
	"silo": {
		ContentType:  "model/mesh",
		Compressible: false,
	},
	"sis": {
		ContentType:  "application/vnd.symbian.install",
		Compressible: false,
	},
	"sisx": {
		ContentType:  "application/vnd.symbian.install",
		Compressible: false,
	},
	"sit": {
		ContentType:  "application/x-stuffit",
		Compressible: false,
	},
	"sitx": {
		ContentType:  "application/x-stuffitx",
		Compressible: false,
	},
	"skd": {
		ContentType:  "application/vnd.koan",
		Compressible: false,
	},
	"skm": {
		ContentType:  "application/vnd.koan",
		Compressible: false,
	},
	"skp": {
		ContentType:  "application/vnd.koan",
		Compressible: false,
	},
	"skt": {
		ContentType:  "application/vnd.koan",
		Compressible: false,
	},
	"sldm": {
		ContentType:  "application/vnd.ms-powerpoint.slide.macroenabled.12",
		Compressible: false,
	},
	"sldx": {
		ContentType:  "application/vnd.openxmlformats-officedocument.presentationml.slide",
		Compressible: false,
	},
	"slim": {
		ContentType:  "text/slim",
		Compressible: false,
	},
	"slm": {
		ContentType:  "text/slim",
		Compressible: false,
	},
	"slt": {
		ContentType:  "application/vnd.epson.salt",
		Compressible: false,
	},
	"sm": {
		ContentType:  "application/vnd.stepmania.stepchart",
		Compressible: false,
	},
	"smf": {
		ContentType:  "application/vnd.stardivision.math",
		Compressible: false,
	},
	"smi": {
		ContentType:  "application/smil+xml",
		Compressible: false,
	},
	"smil": {
		ContentType:  "application/smil+xml",
		Compressible: false,
	},
	"smv": {
		ContentType:  "video/x-smv",
		Compressible: false,
	},
	"smzip": {
		ContentType:  "application/vnd.stepmania.package",
		Compressible: false,
	},
	"snd": {
		ContentType:  "audio/basic",
		Compressible: false,
	},
	"snf": {
		ContentType:  "application/x-font-snf",
		Compressible: false,
	},
	"so": {
		ContentType:  "application/octet-stream",
		Compressible: false,
	},
	"spc": {
		ContentType:  "application/x-pkcs7-certificates",
		Compressible: false,
	},
	"spf": {
		ContentType:  "application/vnd.yamaha.smaf-phrase",
		Compressible: false,
	},
	"spl": {
		ContentType:  "application/x-futuresplash",
		Compressible: false,
	},
	"spot": {
		ContentType:  "text/vnd.in3d.spot",
		Compressible: false,
	},
	"spp": {
		ContentType:  "application/scvp-vp-response",
		Compressible: false,
	},
	"spq": {
		ContentType:  "application/scvp-vp-request",
		Compressible: false,
	},
	"spx": {
		ContentType:  "audio/ogg",
		Compressible: false,
	},
	"sql": {
		ContentType:  "application/x-sql",
		Compressible: false,
	},
	"src": {
		ContentType:  "application/x-wais-source",
		Compressible: false,
	},
	"srt": {
		ContentType:  "application/x-subrip",
		Compressible: false,
	},
	"sru": {
		ContentType:  "application/sru+xml",
		Compressible: false,
	},
	"srx": {
		ContentType:  "application/sparql-results+xml",
		Compressible: false,
	},
	"ssdl": {
		ContentType:  "application/ssdl+xml",
		Compressible: false,
	},
	"sse": {
		ContentType:  "application/vnd.kodak-descriptor",
		Compressible: false,
	},
	"ssf": {
		ContentType:  "application/vnd.epson.ssf",
		Compressible: false,
	},
	"ssml": {
		ContentType:  "application/ssml+xml",
		Compressible: false,
	},
	"st": {
		ContentType:  "application/vnd.sailingtracker.track",
		Compressible: false,
	},
	"stc": {
		ContentType:  "application/vnd.sun.xml.calc.template",
		Compressible: false,
	},
	"std": {
		ContentType:  "application/vnd.sun.xml.draw.template",
		Compressible: false,
	},
	"stf": {
		ContentType:  "application/vnd.wt.stf",
		Compressible: false,
	},
	"sti": {
		ContentType:  "application/vnd.sun.xml.impress.template",
		Compressible: false,
	},
	"stk": {
		ContentType:  "application/hyperstudio",
		Compressible: false,
	},
	"stl": {
		ContentType:  "application/vnd.ms-pki.stl",
		Compressible: false,
	},
	"str": {
		ContentType:  "application/vnd.pg.format",
		Compressible: false,
	},
	"stw": {
		ContentType:  "application/vnd.sun.xml.writer.template",
		Compressible: false,
	},
	"styl": {
		ContentType:  "text/stylus",
		Compressible: false,
	},
	"stylus": {
		ContentType:  "text/stylus",
		Compressible: false,
	},
	"sub": {
		ContentType:  "text/vnd.dvb.subtitle",
		Compressible: false,
	},
	"sus": {
		ContentType:  "application/vnd.sus-calendar",
		Compressible: false,
	},
	"susp": {
		ContentType:  "application/vnd.sus-calendar",
		Compressible: false,
	},
	"sv4cpio": {
		ContentType:  "application/x-sv4cpio",
		Compressible: false,
	},
	"sv4crc": {
		ContentType:  "application/x-sv4crc",
		Compressible: false,
	},
	"svc": {
		ContentType:  "application/vnd.dvb.service",
		Compressible: false,
	},
	"svd": {
		ContentType:  "application/vnd.svd",
		Compressible: false,
	},
	"svg": {
		ContentType:  "image/svg+xml",
		Compressible: false,
	},
	"svgz": {
		ContentType:  "image/svg+xml",
		Compressible: false,
	},
	"swa": {
		ContentType:  "application/x-director",
		Compressible: false,
	},
	"swf": {
		ContentType:  "application/x-shockwave-flash",
		Compressible: false,
	},
	"swi": {
		ContentType:  "application/vnd.aristanetworks.swi",
		Compressible: false,
	},
	"sxc": {
		ContentType:  "application/vnd.sun.xml.calc",
		Compressible: false,
	},
	"sxd": {
		ContentType:  "application/vnd.sun.xml.draw",
		Compressible: false,
	},
	"sxg": {
		ContentType:  "application/vnd.sun.xml.writer.global",
		Compressible: false,
	},
	"sxi": {
		ContentType:  "application/vnd.sun.xml.impress",
		Compressible: false,
	},
	"sxm": {
		ContentType:  "application/vnd.sun.xml.math",
		Compressible: false,
	},
	"sxw": {
		ContentType:  "application/vnd.sun.xml.writer",
		Compressible: false,
	},
	"t": {
		ContentType:  "text/troff",
		Compressible: false,
	},
	"t3": {
		ContentType:  "application/x-t3vm-image",
		Compressible: false,
	},
	"t38": {
		ContentType:  "image/t38",
		Compressible: false,
	},
	"taglet": {
		ContentType:  "application/vnd.mynfc",
		Compressible: false,
	},
	"tao": {
		ContentType:  "application/vnd.tao.intent-module-archive",
		Compressible: false,
	},
	"tap": {
		ContentType:  "image/vnd.tencent.tap",
		Compressible: false,
	},
	"tar": {
		ContentType:  "application/x-tar",
		Compressible: false,
	},
	"tcap": {
		ContentType:  "application/vnd.3gpp2.tcap",
		Compressible: false,
	},
	"tcl": {
		ContentType:  "application/x-tcl",
		Compressible: false,
	},
	"teacher": {
		ContentType:  "application/vnd.smart.teacher",
		Compressible: false,
	},
	"tei": {
		ContentType:  "application/tei+xml",
		Compressible: false,
	},
	"teicorpus": {
		ContentType:  "application/tei+xml",
		Compressible: false,
	},
	"tex": {
		ContentType:  "application/x-tex",
		Compressible: false,
	},
	"texi": {
		ContentType:  "application/x-texinfo",
		Compressible: false,
	},
	"texinfo": {
		ContentType:  "application/x-texinfo",
		Compressible: false,
	},
	"text": {
		ContentType:  "text/plain",
		Compressible: false,
	},
	"tfi": {
		ContentType:  "application/thraud+xml",
		Compressible: false,
	},
	"tfm": {
		ContentType:  "application/x-tex-tfm",
		Compressible: false,
	},
	"tfx": {
		ContentType:  "image/tiff-fx",
		Compressible: false,
	},
	"tga": {
		ContentType:  "image/x-tga",
		Compressible: false,
	},
	"thmx": {
		ContentType:  "application/vnd.ms-officetheme",
		Compressible: false,
	},
	"tif": {
		ContentType:  "image/tiff",
		Compressible: false,
	},
	"tiff": {
		ContentType:  "image/tiff",
		Compressible: false,
	},
	"tk": {
		ContentType:  "application/x-tcl",
		Compressible: false,
	},
	"tmo": {
		ContentType:  "application/vnd.tmobile-livetv",
		Compressible: false,
	},
	"torrent": {
		ContentType:  "application/x-bittorrent",
		Compressible: false,
	},
	"tpl": {
		ContentType:  "application/vnd.groove-tool-template",
		Compressible: false,
	},
	"tpt": {
		ContentType:  "application/vnd.trid.tpt",
		Compressible: false,
	},
	"tr": {
		ContentType:  "text/troff",
		Compressible: false,
	},
	"tra": {
		ContentType:  "application/vnd.trueapp",
		Compressible: false,
	},
	"trm": {
		ContentType:  "application/x-msterminal",
		Compressible: false,
	},
	"ts": {
		ContentType:  "video/mp2t",
		Compressible: false,
	},
	"tsd": {
		ContentType:  "application/timestamped-data",
		Compressible: false,
	},
	"tsv": {
		ContentType:  "text/tab-separated-values",
		Compressible: false,
	},
	"ttc": {
		ContentType:  "font/collection",
		Compressible: false,
	},
	"ttf": {
		ContentType:  "font/ttf",
		Compressible: false,
	},
	"ttl": {
		ContentType:  "text/turtle",
		Compressible: false,
	},
	"twd": {
		ContentType:  "application/vnd.simtech-mindmapper",
		Compressible: false,
	},
	"twds": {
		ContentType:  "application/vnd.simtech-mindmapper",
		Compressible: false,
	},
	"txd": {
		ContentType:  "application/vnd.genomatix.tuxedo",
		Compressible: false,
	},
	"txf": {
		ContentType:  "application/vnd.mobius.txf",
		Compressible: false,
	},
	"txt": {
		ContentType:  "text/plain",
		Compressible: false,
	},
	"u32": {
		ContentType:  "application/x-authorware-bin",
		Compressible: false,
	},
	"u8dsn": {
		ContentType:  "message/global-delivery-status",
		Compressible: false,
	},
	"u8hdr": {
		ContentType:  "message/global-headers",
		Compressible: false,
	},
	"u8mdn": {
		ContentType:  "message/global-disposition-notification",
		Compressible: false,
	},
	"u8msg": {
		ContentType:  "message/global",
		Compressible: false,
	},
	"udeb": {
		ContentType:  "application/x-debian-package",
		Compressible: false,
	},
	"ufd": {
		ContentType:  "application/vnd.ufdl",
		Compressible: false,
	},
	"ufdl": {
		ContentType:  "application/vnd.ufdl",
		Compressible: false,
	},
	"ulx": {
		ContentType:  "application/x-glulx",
		Compressible: false,
	},
	"umj": {
		ContentType:  "application/vnd.umajin",
		Compressible: false,
	},
	"unityweb": {
		ContentType:  "application/vnd.unity",
		Compressible: false,
	},
	"uoml": {
		ContentType:  "application/vnd.uoml+xml",
		Compressible: false,
	},
	"uri": {
		ContentType:  "text/uri-list",
		Compressible: false,
	},
	"uris": {
		ContentType:  "text/uri-list",
		Compressible: false,
	},
	"urls": {
		ContentType:  "text/uri-list",
		Compressible: false,
	},
	"ustar": {
		ContentType:  "application/x-ustar",
		Compressible: false,
	},
	"utz": {
		ContentType:  "application/vnd.uiq.theme",
		Compressible: false,
	},
	"uu": {
		ContentType:  "text/x-uuencode",
		Compressible: false,
	},
	"uva": {
		ContentType:  "audio/vnd.dece.audio",
		Compressible: false,
	},
	"uvd": {
		ContentType:  "application/vnd.dece.data",
		Compressible: false,
	},
	"uvf": {
		ContentType:  "application/vnd.dece.data",
		Compressible: false,
	},
	"uvg": {
		ContentType:  "image/vnd.dece.graphic",
		Compressible: false,
	},
	"uvh": {
		ContentType:  "video/vnd.dece.hd",
		Compressible: false,
	},
	"uvi": {
		ContentType:  "image/vnd.dece.graphic",
		Compressible: false,
	},
	"uvm": {
		ContentType:  "video/vnd.dece.mobile",
		Compressible: false,
	},
	"uvp": {
		ContentType:  "video/vnd.dece.pd",
		Compressible: false,
	},
	"uvs": {
		ContentType:  "video/vnd.dece.sd",
		Compressible: false,
	},
	"uvt": {
		ContentType:  "application/vnd.dece.ttml+xml",
		Compressible: false,
	},
	"uvu": {
		ContentType:  "video/vnd.uvvu.mp4",
		Compressible: false,
	},
	"uvv": {
		ContentType:  "video/vnd.dece.video",
		Compressible: false,
	},
	"uvva": {
		ContentType:  "audio/vnd.dece.audio",
		Compressible: false,
	},
	"uvvd": {
		ContentType:  "application/vnd.dece.data",
		Compressible: false,
	},
	"uvvf": {
		ContentType:  "application/vnd.dece.data",
		Compressible: false,
	},
	"uvvg": {
		ContentType:  "image/vnd.dece.graphic",
		Compressible: false,
	},
	"uvvh": {
		ContentType:  "video/vnd.dece.hd",
		Compressible: false,
	},
	"uvvi": {
		ContentType:  "image/vnd.dece.graphic",
		Compressible: false,
	},
	"uvvm": {
		ContentType:  "video/vnd.dece.mobile",
		Compressible: false,
	},
	"uvvp": {
		ContentType:  "video/vnd.dece.pd",
		Compressible: false,
	},
	"uvvs": {
		ContentType:  "video/vnd.dece.sd",
		Compressible: false,
	},
	"uvvt": {
		ContentType:  "application/vnd.dece.ttml+xml",
		Compressible: false,
	},
	"uvvu": {
		ContentType:  "video/vnd.uvvu.mp4",
		Compressible: false,
	},
	"uvvv": {
		ContentType:  "video/vnd.dece.video",
		Compressible: false,
	},
	"uvvx": {
		ContentType:  "application/vnd.dece.unspecified",
		Compressible: false,
	},
	"uvvz": {
		ContentType:  "application/vnd.dece.zip",
		Compressible: false,
	},
	"uvx": {
		ContentType:  "application/vnd.dece.unspecified",
		Compressible: false,
	},
	"uvz": {
		ContentType:  "application/vnd.dece.zip",
		Compressible: false,
	},
	"vbox": {
		ContentType:  "application/x-virtualbox-vbox",
		Compressible: false,
	},
	"vbox-extpack": {
		ContentType:  "application/x-virtualbox-vbox-extpack",
		Compressible: false,
	},
	"vcard": {
		ContentType:  "text/vcard",
		Compressible: false,
	},
	"vcd": {
		ContentType:  "application/x-cdlink",
		Compressible: false,
	},
	"vcf": {
		ContentType:  "text/x-vcard",
		Compressible: false,
	},
	"vcg": {
		ContentType:  "application/vnd.groove-vcard",
		Compressible: false,
	},
	"vcs": {
		ContentType:  "text/x-vcalendar",
		Compressible: false,
	},
	"vcx": {
		ContentType:  "application/vnd.vcx",
		Compressible: false,
	},
	"vdi": {
		ContentType:  "application/x-virtualbox-vdi",
		Compressible: false,
	},
	"vhd": {
		ContentType:  "application/x-virtualbox-vhd",
		Compressible: false,
	},
	"vis": {
		ContentType:  "application/vnd.visionary",
		Compressible: false,
	},
	"viv": {
		ContentType:  "video/vnd.vivo",
		Compressible: false,
	},
	"vmdk": {
		ContentType:  "application/x-virtualbox-vmdk",
		Compressible: false,
	},
	"vob": {
		ContentType:  "video/x-ms-vob",
		Compressible: false,
	},
	"vor": {
		ContentType:  "application/vnd.stardivision.writer",
		Compressible: false,
	},
	"vox": {
		ContentType:  "application/x-authorware-bin",
		Compressible: false,
	},
	"vrml": {
		ContentType:  "model/vrml",
		Compressible: false,
	},
	"vsd": {
		ContentType:  "application/vnd.visio",
		Compressible: false,
	},
	"vsf": {
		ContentType:  "application/vnd.vsf",
		Compressible: false,
	},
	"vss": {
		ContentType:  "application/vnd.visio",
		Compressible: false,
	},
	"vst": {
		ContentType:  "application/vnd.visio",
		Compressible: false,
	},
	"vsw": {
		ContentType:  "application/vnd.visio",
		Compressible: false,
	},
	"vtf": {
		ContentType:  "image/vnd.valve.source.texture",
		Compressible: false,
	},
	"vtt": {
		ContentType:  "text/vtt",
		Compressible: false,
	},
	"vtu": {
		ContentType:  "model/vnd.vtu",
		Compressible: false,
	},
	"vxml": {
		ContentType:  "application/voicexml+xml",
		Compressible: false,
	},
	"w3d": {
		ContentType:  "application/x-director",
		Compressible: false,
	},
	"wad": {
		ContentType:  "application/x-doom",
		Compressible: false,
	},
	"wadl": {
		ContentType:  "application/vnd.sun.wadl+xml",
		Compressible: false,
	},
	"war": {
		ContentType:  "application/java-archive",
		Compressible: false,
	},
	"wasm": {
		ContentType:  "application/wasm",
		Compressible: false,
	},
	"wav": {
		ContentType:  "audio/x-wav",
		Compressible: false,
	},
	"wax": {
		ContentType:  "audio/x-ms-wax",
		Compressible: false,
	},
	"wbmp": {
		ContentType:  "image/vnd.wap.wbmp",
		Compressible: false,
	},
	"wbs": {
		ContentType:  "application/vnd.criticaltools.wbs+xml",
		Compressible: false,
	},
	"wbxml": {
		ContentType:  "application/vnd.wap.wbxml",
		Compressible: false,
	},
	"wcm": {
		ContentType:  "application/vnd.ms-works",
		Compressible: false,
	},
	"wdb": {
		ContentType:  "application/vnd.ms-works",
		Compressible: false,
	},
	"wdp": {
		ContentType:  "image/vnd.ms-photo",
		Compressible: false,
	},
	"weba": {
		ContentType:  "audio/webm",
		Compressible: false,
	},
	"webapp": {
		ContentType:  "application/x-web-app-manifest+json",
		Compressible: false,
	},
	"webm": {
		ContentType:  "video/webm",
		Compressible: false,
	},
	"webmanifest": {
		ContentType:  "application/manifest+json",
		Compressible: false,
	},
	"webp": {
		ContentType:  "image/webp",
		Compressible: false,
	},
	"wg": {
		ContentType:  "application/vnd.pmi.widget",
		Compressible: false,
	},
	"wgt": {
		ContentType:  "application/widget",
		Compressible: false,
	},
	"wks": {
		ContentType:  "application/vnd.ms-works",
		Compressible: false,
	},
	"wm": {
		ContentType:  "video/x-ms-wm",
		Compressible: false,
	},
	"wma": {
		ContentType:  "audio/x-ms-wma",
		Compressible: false,
	},
	"wmd": {
		ContentType:  "application/x-ms-wmd",
		Compressible: false,
	},
	"wmf": {
		ContentType:  "image/wmf",
		Compressible: false,
	},
	"wml": {
		ContentType:  "text/vnd.wap.wml",
		Compressible: false,
	},
	"wmlc": {
		ContentType:  "application/vnd.wap.wmlc",
		Compressible: false,
	},
	"wmls": {
		ContentType:  "text/vnd.wap.wmlscript",
		Compressible: false,
	},
	"wmlsc": {
		ContentType:  "application/vnd.wap.wmlscriptc",
		Compressible: false,
	},
	"wmv": {
		ContentType:  "video/x-ms-wmv",
		Compressible: false,
	},
	"wmx": {
		ContentType:  "video/x-ms-wmx",
		Compressible: false,
	},
	"wmz": {
		ContentType:  "application/x-msmetafile",
		Compressible: false,
	},
	"woff": {
		ContentType:  "font/woff",
		Compressible: false,
	},
	"woff2": {
		ContentType:  "font/woff2",
		Compressible: false,
	},
	"wpd": {
		ContentType:  "application/vnd.wordperfect",
		Compressible: false,
	},
	"wpl": {
		ContentType:  "application/vnd.ms-wpl",
		Compressible: false,
	},
	"wps": {
		ContentType:  "application/vnd.ms-works",
		Compressible: false,
	},
	"wqd": {
		ContentType:  "application/vnd.wqd",
		Compressible: false,
	},
	"wri": {
		ContentType:  "application/x-mswrite",
		Compressible: false,
	},
	"wrl": {
		ContentType:  "model/vrml",
		Compressible: false,
	},
	"wsc": {
		ContentType:  "message/vnd.wfa.wsc",
		Compressible: false,
	},
	"wsdl": {
		ContentType:  "application/wsdl+xml",
		Compressible: false,
	},
	"wspolicy": {
		ContentType:  "application/wspolicy+xml",
		Compressible: false,
	},
	"wtb": {
		ContentType:  "application/vnd.webturbo",
		Compressible: false,
	},
	"wvx": {
		ContentType:  "video/x-ms-wvx",
		Compressible: false,
	},
	"x32": {
		ContentType:  "application/x-authorware-bin",
		Compressible: false,
	},
	"x3d": {
		ContentType:  "model/x3d+xml",
		Compressible: false,
	},
	"x3db": {
		ContentType:  "model/x3d+binary",
		Compressible: false,
	},
	"x3dbz": {
		ContentType:  "model/x3d+binary",
		Compressible: false,
	},
	"x3dv": {
		ContentType:  "model/x3d+vrml",
		Compressible: false,
	},
	"x3dvz": {
		ContentType:  "model/x3d+vrml",
		Compressible: false,
	},
	"x3dz": {
		ContentType:  "model/x3d+xml",
		Compressible: false,
	},
	"xaml": {
		ContentType:  "application/xaml+xml",
		Compressible: false,
	},
	"xap": {
		ContentType:  "application/x-silverlight-app",
		Compressible: false,
	},
	"xar": {
		ContentType:  "application/vnd.xara",
		Compressible: false,
	},
	"xbap": {
		ContentType:  "application/x-ms-xbap",
		Compressible: false,
	},
	"xbd": {
		ContentType:  "application/vnd.fujixerox.docuworks.binder",
		Compressible: false,
	},
	"xbm": {
		ContentType:  "image/x-xbitmap",
		Compressible: false,
	},
	"xdf": {
		ContentType:  "application/xcap-diff+xml",
		Compressible: false,
	},
	"xdm": {
		ContentType:  "application/vnd.syncml.dm+xml",
		Compressible: false,
	},
	"xdp": {
		ContentType:  "application/vnd.adobe.xdp+xml",
		Compressible: false,
	},
	"xdssc": {
		ContentType:  "application/dssc+xml",
		Compressible: false,
	},
	"xdw": {
		ContentType:  "application/vnd.fujixerox.docuworks",
		Compressible: false,
	},
	"xenc": {
		ContentType:  "application/xenc+xml",
		Compressible: false,
	},
	"xer": {
		ContentType:  "application/patch-ops-error+xml",
		Compressible: false,
	},
	"xfdf": {
		ContentType:  "application/vnd.adobe.xfdf",
		Compressible: false,
	},
	"xfdl": {
		ContentType:  "application/vnd.xfdl",
		Compressible: false,
	},
	"xht": {
		ContentType:  "application/xhtml+xml",
		Compressible: false,
	},
	"xhtml": {
		ContentType:  "application/xhtml+xml",
		Compressible: false,
	},
	"xhvml": {
		ContentType:  "application/xv+xml",
		Compressible: false,
	},
	"xif": {
		ContentType:  "image/vnd.xiff",
		Compressible: false,
	},
	"xla": {
		ContentType:  "application/vnd.ms-excel",
		Compressible: false,
	},
	"xlam": {
		ContentType:  "application/vnd.ms-excel.addin.macroenabled.12",
		Compressible: false,
	},
	"xlc": {
		ContentType:  "application/vnd.ms-excel",
		Compressible: false,
	},
	"xlf": {
		ContentType:  "application/x-xliff+xml",
		Compressible: false,
	},
	"xlm": {
		ContentType:  "application/vnd.ms-excel",
		Compressible: false,
	},
	"xls": {
		ContentType:  "application/vnd.ms-excel",
		Compressible: false,
	},
	"xlsb": {
		ContentType:  "application/vnd.ms-excel.sheet.binary.macroenabled.12",
		Compressible: false,
	},
	"xlsm": {
		ContentType:  "application/vnd.ms-excel.sheet.macroenabled.12",
		Compressible: false,
	},
	"xlsx": {
		ContentType:  "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
		Compressible: false,
	},
	"xlt": {
		ContentType:  "application/vnd.ms-excel",
		Compressible: false,
	},
	"xltm": {
		ContentType:  "application/vnd.ms-excel.template.macroenabled.12",
		Compressible: false,
	},
	"xltx": {
		ContentType:  "application/vnd.openxmlformats-officedocument.spreadsheetml.template",
		Compressible: false,
	},
	"xlw": {
		ContentType:  "application/vnd.ms-excel",
		Compressible: false,
	},
	"xm": {
		ContentType:  "audio/xm",
		Compressible: false,
	},
	"xml": {
		ContentType:  "text/xml",
		Compressible: false,
	},
	"xo": {
		ContentType:  "application/vnd.olpc-sugar",
		Compressible: false,
	},
	"xop": {
		ContentType:  "application/xop+xml",
		Compressible: false,
	},
	"xpi": {
		ContentType:  "application/x-xpinstall",
		Compressible: false,
	},
	"xpl": {
		ContentType:  "application/xproc+xml",
		Compressible: false,
	},
	"xpm": {
		ContentType:  "image/x-xpixmap",
		Compressible: false,
	},
	"xpr": {
		ContentType:  "application/vnd.is-xpr",
		Compressible: false,
	},
	"xps": {
		ContentType:  "application/vnd.ms-xpsdocument",
		Compressible: false,
	},
	"xpw": {
		ContentType:  "application/vnd.intercon.formnet",
		Compressible: false,
	},
	"xpx": {
		ContentType:  "application/vnd.intercon.formnet",
		Compressible: false,
	},
	"xsd": {
		ContentType:  "application/xml",
		Compressible: false,
	},
	"xsl": {
		ContentType:  "application/xml",
		Compressible: false,
	},
	"xslt": {
		ContentType:  "application/xslt+xml",
		Compressible: false,
	},
	"xsm": {
		ContentType:  "application/vnd.syncml+xml",
		Compressible: false,
	},
	"xspf": {
		ContentType:  "application/xspf+xml",
		Compressible: false,
	},
	"xul": {
		ContentType:  "application/vnd.mozilla.xul+xml",
		Compressible: false,
	},
	"xvm": {
		ContentType:  "application/xv+xml",
		Compressible: false,
	},
	"xvml": {
		ContentType:  "application/xv+xml",
		Compressible: false,
	},
	"xwd": {
		ContentType:  "image/x-xwindowdump",
		Compressible: false,
	},
	"xyz": {
		ContentType:  "chemical/x-xyz",
		Compressible: false,
	},
	"xz": {
		ContentType:  "application/x-xz",
		Compressible: false,
	},
	"yaml": {
		ContentType:  "text/yaml",
		Compressible: false,
	},
	"yang": {
		ContentType:  "application/yang",
		Compressible: false,
	},
	"yin": {
		ContentType:  "application/yin+xml",
		Compressible: false,
	},
	"yml": {
		ContentType:  "text/yaml",
		Compressible: false,
	},
	"ymp": {
		ContentType:  "text/x-suse-ymp",
		Compressible: false,
	},
	"z1": {
		ContentType:  "application/x-zmachine",
		Compressible: false,
	},
	"z2": {
		ContentType:  "application/x-zmachine",
		Compressible: false,
	},
	"z3": {
		ContentType:  "application/x-zmachine",
		Compressible: false,
	},
	"z4": {
		ContentType:  "application/x-zmachine",
		Compressible: false,
	},
	"z5": {
		ContentType:  "application/x-zmachine",
		Compressible: false,
	},
	"z6": {
		ContentType:  "application/x-zmachine",
		Compressible: false,
	},
	"z7": {
		ContentType:  "application/x-zmachine",
		Compressible: false,
	},
	"z8": {
		ContentType:  "application/x-zmachine",
		Compressible: false,
	},
	"zaz": {
		ContentType:  "application/vnd.zzazz.deck+xml",
		Compressible: false,
	},
	"zip": {
		ContentType:  "application/zip",
		Compressible: false,
	},
	"zir": {
		ContentType:  "application/vnd.zul",
		Compressible: false,
	},
	"zirz": {
		ContentType:  "application/vnd.zul",
		Compressible: false,
	},
	"zmm": {
		ContentType:  "application/vnd.handheld-entertainment+xml",
		Compressible: false,
	},
}
