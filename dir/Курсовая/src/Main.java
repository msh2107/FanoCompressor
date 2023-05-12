import java.io.File;
import java.io.RandomAccessFile;
import java.util.HexFormat;

class Main {

    public static void main(String[] args) throws Exception {
        if (Utils.processInput(args)) {
            DES.runTests();
        } else {
            DES.setKeys(Utils.passwordToKey());

            for (File file : Utils.getFiles()) {
                DES.setInitVector(Utils.getInitVector());

                RandomAccessFile rwFile = new RandomAccessFile(file, "rw");
                byte[] buffer = new byte[8];

                int read = rwFile.read(buffer);
                while (read != -1) {
                    rwFile.seek(rwFile.getFilePointer() - read);

                    String hexString = DES.encryptDecryptOFB(HexFormat.of().formatHex(buffer));
                    buffer = HexFormat.of().parseHex(hexString);
                    for (int j = 0; j < read; j++) {
                        rwFile.writeByte(buffer[j]);
                    }

                    read = rwFile.read(buffer);
                }

                rwFile.close();
            }
        }
    }
}
